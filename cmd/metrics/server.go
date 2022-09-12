/*
Copyright The Kubeall Authors.

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
package metrics

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/kube-all/metrics-collector/cmd/metrics/options"
	"github.com/kube-all/metrics-collector/pkg/embeds"
	"github.com/kube-all/metrics-collector/pkg/global"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

func KubeMetricsGenerateCommand() *cobra.Command {
	s := options.NewKubeMetricsGenerateOptions()
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "generate default metrics",
		Version: global.Version,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				s.Path = args[0]
			}
			if err := Run(s); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	s.AddFlag(cmd.Flags())
	return cmd
}

func Run(s *options.KubeMetricsGenerateOptions) (err error) {
	if len(s.Path) == 0 {
		err = errors.New("you must specify generate path with parameter -path or -p")
		return err
	}
	_, err = os.Stat(s.Path)
	if os.IsNotExist(err) {
		os.MkdirAll(path.Join(s.Path, "metrics"), os.ModePerm)
	}
	// generate metrics config
	files, _ := embeds.MetricsYamlFiles.ReadDir("metrics")
	for _, f := range files {
		if data, e := embeds.MetricsYamlFiles.ReadFile(path.Join("metrics", f.Name())); e == nil {
			fp := path.Join(s.Path, "metrics", f.Name())
			if e1 := ioutil.WriteFile(fp, data, os.ModePerm); e1 == nil {
				klog.Infof("generate file: %s success", fp)
			} else {
				klog.Errorf("generate file: %s failed, err: %s", f.Name(), e1.Error())
			}
		}
	}
	// generate clusters information config
	if data, e := embeds.DemoClusters.ReadFile("clusters.yaml"); e == nil {
		fp := path.Join(s.Path, "clusters.yaml")
		if e1 := ioutil.WriteFile(fp, data, os.ModePerm); e1 == nil {
			klog.Infof("generate file: %s success", fp)
		} else {
			klog.Errorf("generate file: clusters.yaml failed, err: %s", e1.Error())
		}
	}
	// generate app config
	global.ApplicationConfig.Collectors = []string{"cluster", "namespace", "node", "service", "pod"}
	global.ApplicationConfig.ClusterPath = path.Join(s.Path, "clusters.yaml")
	global.ApplicationConfig.MetricsPath = path.Join(s.Path, "metrics")
	if d, e := yaml.Marshal(global.ApplicationConfig); e == nil {

		if e1 := ioutil.WriteFile(path.Join(s.Path, "config.yaml"), d, os.ModePerm); e1 == nil {
			klog.Infof("generate metrics-collector file: %s success", path.Join(s.Path, "config.yaml"))
		} else {
			klog.Infof("generate metrics-collector file: %s failed, err: %s",
				path.Join(s.Path, "config.yaml"), e1.Error())
		}
	}
	return nil
}
