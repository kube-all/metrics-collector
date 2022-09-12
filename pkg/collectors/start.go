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
	"k8s.io/klog/v2"
	"sync"
)

func collector(code string) {
	me := st.GetMetrics(code)
	if len(me.MetricsPromQLs) == 0 {
		klog.Warning("no info to start collector")
		return
	}
	klog.V(6).Infof("start collector job for code: %s, name: %s, cron expression: %s, metrics number: %d",
		code, me.Name, me.CronExpression, len(me.MetricsPromQLs))
	var wg sync.WaitGroup
	for _, cluster := range st.GetClusters() {
		wg.Add(1)
		go func(c Cluster, m MetricsCategory) {
			collectorJob(c, m)
			wg.Done()
		}(cluster, me)
	}
	wg.Wait()
	klog.V(6).Infof("end collector job for code: %s, name: %s, cron expression: %s, metrics number: %d",
		code, me.Name, me.CronExpression, len(me.MetricsPromQLs))
}
