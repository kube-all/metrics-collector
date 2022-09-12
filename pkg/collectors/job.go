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
	"github.com/robfig/cron/v3"
	"k8s.io/klog/v2"
	"sync"
)

type CronJob struct {
	lock               sync.RWMutex
	Cron               *cron.Cron
	Entries            map[string]cron.EntryID
	CronExpressions    map[string]string
	NewCronExpressions map[string]string
	Started            bool
}

var jobOnce sync.Once

var job *CronJob

func GetCronJobInstance() *CronJob {
	if job == nil {
		jobOnce.Do(func() {
			job = &CronJob{
				Cron:               cron.New(),
				Entries:            make(map[string]cron.EntryID),
				CronExpressions:    make(map[string]string),
				NewCronExpressions: make(map[string]string),
			}
		})
	}
	return job
}
func (job *CronJob) Add(jobName string, me MetricsCategory) {
	job.lock.Lock()
	defer job.lock.Unlock()
	job.NewCronExpressions[jobName] = me.CronExpression
}

func (job *CronJob) Start() {
	job.lock.Lock()
	defer job.lock.Unlock()
	var okJobs []string
	for jb, exp := range job.NewCronExpressions {
		var (
			jobName    = jb
			expression = exp
		)
		//任务没有启动
		if ex, ok := job.CronExpressions[jobName]; !ok {
			if entId, err := job.Cron.AddFunc(expression, func() { collector(jobName) }); err == nil {
				job.Entries[jobName] = entId
				job.CronExpressions[jobName] = expression
				okJobs = append(okJobs, jobName)
				klog.V(6).Infof("add cron collector job for: %s success, cron expression: %s", jobName, expression)
			} else {
				klog.Infof("add cron collector job for: %s failed, err: %s", jobName, err.Error())
			}
		} else if ex != expression {
			if id, ok := job.Entries[jobName]; ok {
				job.Cron.Remove(id)
				delete(job.CronExpressions, jobName)
				klog.Infof("delete cron collector job for: %s success", jobName)
			}
			if entId, err := job.Cron.AddFunc(expression, func() { collector(jobName) }); err == nil {
				job.Entries[jobName] = entId
				job.CronExpressions[jobName] = expression
				okJobs = append(okJobs, jobName)
				klog.Infof("add cron collector job for: %s success, cron expression: %s", jobName, expression)
			} else {
				klog.Infof("add cron collector job for: %s failed, err: %s", jobName, err.Error())
			}
		} else {
			okJobs = append(okJobs, jobName)
		}
	}
	for _, j := range okJobs {
		if _, exist := job.NewCronExpressions[j]; exist {
			delete(job.NewCronExpressions, j)
		}
	}
	if !job.Started {
		job.Cron.Start()
		job.Started = true
	}
}
