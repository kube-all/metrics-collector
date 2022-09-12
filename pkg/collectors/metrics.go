/*
Copyright 2022 The kubeall.com Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collectors

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/kube-all/metrics-collector/pkg/utils"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

var stOnce sync.Once

var st *CollectorStorage

type Cluster struct {
	Name       string     `json:"name,omitempty" yaml:"name" description:"集群名称" validate:"required"`
	Code       string     `json:"code,omitempty" yaml:"code" description:"集群编码"`
	Prometheus Prometheus `json:"prometheus,omitempty" yaml:"prometheus"`
}
type Prometheus struct {
	Address  string `json:"address,omitempty" yaml:"address"`
	Token    string `json:"token,omitempty" yaml:"token"`
	Username string `json:"username,omitempty" yaml:"username"`
	Password string `json:"password,omitempty" yaml:"password"`
}
type CollectorStorage struct {
	lock           sync.RWMutex
	Clusters       map[string]Cluster         `json:"clusters,omitempty" yaml:"clusters"`
	Metrics        map[string]MetricsCategory `json:"metrics,omitempty" yaml:"metrics"`
	Collectors     []string                   `json:"collectors,omitempty" yaml:"collectors"`
	CronExpression map[string]string          `json:"-" yaml:"-"`
	MetricsPath    string                     `json:"-" yaml:"-"`
	ClusterPath    string                     `json:"-" yaml:"-"`
}

type MetricsCategory struct {
	Name           string           `yaml:"name" json:"name" validate:"required" example:"average"`                         // 类型名称
	Code           string           `json:"code,omitempty" yaml:"code" validate:"oneof=cluster node namespace service pod"` // 类型编码
	Description    string           `json:"description" yaml:"description" example:"one hour average value "`               // 说明
	CronExpression string           `json:"cronExpression" validate:"required" yaml:"cronExpression" example:"@every 2h"`   // 定时任务表达式
	MetricsPromQLs []*MetricsPromQL `json:"promQls" yaml:"promQls"`                                                         // promql语句，用于首次从yaml导入
	ExtractLabels  []*ExtractLabel  `json:"extractLabels,omitempty" yaml:"extractLabels"`
}
type ExtractLabel struct {
	Label   string `json:"label,omitempty" yaml:"label"`     // 标签名称
	DBField string `json:"dbField,omitempty" yaml:"dbField"` // 映射到数据库字段
}
type MetricsPromQL struct {
	Name        string `json:"name,omitempty" yaml:"name" validate:"required"`                         // promql语句名称，功能
	Description string `json:"description,omitempty" yaml:"description"`                               // 指标描述
	Code        string `json:"code" validate:"required,check"`                                         // promql语句编码
	Query       string `json:"query"`                                                                  // prometheus查询语句
	ValueType   string `json:"valueType,omitempty" yaml:"valueType" validate:"oneof=int string float"` // 指标结果类型
	Unit        string `json:"unit,omitempty" yaml:"unit" validate:"oneof=M C ''"`
}

func (m *MetricsCategory) DefaultValue() {
	for i, ins := range m.ExtractLabels {
		if len(ins.DBField) == 0 {
			ins.DBField = ins.Label
			m.ExtractLabels[i] = ins
		}
	}
}
func (m *MetricsCategory) Validator() (errs []error) {
	labels := make(map[string]string)
	for _, ins := range m.ExtractLabels {
		labels[ins.Label] = ins.Label
	}
	if len(labels) != len(m.ExtractLabels) {
		errs = append(errs, fmt.Errorf("duplicate ExtractLabels exists in : %s", m.Name))
	}
	codes := make(map[string]string)
	for _, ins := range m.MetricsPromQLs {
		codes[ins.Code] = ins.Code
	}
	if len(codes) != len(m.MetricsPromQLs) {
		errs = append(errs, fmt.Errorf("duplicate MetricsPromQLs exists in : %s", m.Name))
	}
	return errs
}

func GetCollectorStorageInstance(clusterPath, metricsPath string, collectors []string) *CollectorStorage {
	if st == nil {
		stOnce.Do(func() {
			st = &CollectorStorage{
				Collectors:     collectors,
				MetricsPath:    metricsPath,
				ClusterPath:    clusterPath,
				Clusters:       make(map[string]Cluster),
				Metrics:        make(map[string]MetricsCategory),
				CronExpression: make(map[string]string),
			}
		})
	}
	return st
}
func (c *CollectorStorage) Load() {
	c.lock.Lock()
	defer c.lock.Unlock()
	clusters := make(map[string]Cluster)

	if len(c.MetricsPath) > 0 {
		filepath.Walk(c.MetricsPath, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				klog.Warningf("load data from path: %s failed, err: %s", p, err.Error())
			}
			if info != nil && info.IsDir() {
				return nil
			}
			relPath, err := filepath.Rel(c.MetricsPath, p)
			if err != nil {
				klog.Warningf("file path: %s rel: %s, err: %s ", p, c.MetricsPath, err.Error())
			}

			if strings.HasSuffix(relPath, ".yaml") || strings.HasSuffix(relPath, ".yml") {
				if data, e := ioutil.ReadFile(p); e == nil {
					var me MetricsCategory
					if e1 := yaml.Unmarshal(data, &me); e1 == nil {
						if utils.StringKeyInArray(me.Code, c.Collectors) {
							c.Metrics[me.Code] = me
							GetCronJobInstance().Add(me.Code, me)
							klog.Infof("load metrics file: %s success, Name: %s, Code: %s,CronExpression: %s",
								relPath, me.Name, me.Code, me.CronExpression)
						} else {
							klog.Warningf("file: %s will not collect it's metrics, you need add [%s] to config file collectors",
								relPath, me.Code)
						}

					} else {
						klog.Errorf("yaml Unmarshal file: %s failed, err: %s", p, e1.Error())
					}
				} else {
					klog.Errorf("load metrics file: %s failed, err: %s", relPath, e.Error())
				}
			}
			return nil
		})
	}
	if len(c.ClusterPath) > 0 {
		if data, e := ioutil.ReadFile(c.ClusterPath); e == nil {
			var cs []Cluster
			if e1 := yaml.Unmarshal(data, &cs); e1 == nil {
				for _, cluster := range cs {
					c.Clusters[cluster.Code] = cluster
					clusters[cluster.Code] = cluster
					klog.Infof("load cluster: %s success", cluster.Code)
				}
			} else {
				klog.Errorf("yaml Unmarshal file: %s failed, err: %s", c.ClusterPath, e1.Error())
			}
		} else {
			klog.Errorf("read file: %s failed, err: %s", c.ClusterPath, e.Error())
		}
	}
	GetCronJobInstance().Start()
	klog.Info("load file success")
}
func (c *CollectorStorage) Watcher(collectors []string) {
	// 第一次先加载配置
	c.lock.Lock()
	c.Collectors = collectors
	c.lock.Unlock()
	c.Load()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.Errorf("NewWatcher failed: %s", err.Error())
		return
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if strings.HasSuffix(event.Name, ".yaml") || strings.HasSuffix(event.Name, ".yml") {
					c.Load()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				klog.Info("error:", err)
			}
		}
	}()

	err = watcher.Add(c.MetricsPath)
	if err != nil {
		klog.Errorf("Add %s failed, err: %s", c.MetricsPath, err.Error())
	} else {
		klog.Infof("dir: %s watch success", c.MetricsPath)
	}
	err = watcher.Add(c.ClusterPath)
	if err != nil {
		klog.Errorf("Add %s failed, err: %s", c.ClusterPath, err.Error())
	} else {
		klog.Infof("file: %s watch success", c.ClusterPath)

	}
	<-done

}
func (c *CollectorStorage) GetMetrics(code string) (me MetricsCategory) {
	c.lock.Lock()
	defer c.lock.Unlock()
	me, _ = c.Metrics[code]
	return
}
func (c *CollectorStorage) GetClusters() (clusters map[string]Cluster) {
	c.lock.Lock()
	defer c.lock.Unlock()
	clusters = c.Clusters
	return
}
