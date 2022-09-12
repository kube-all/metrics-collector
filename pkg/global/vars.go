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

package global

import "github.com/kube-all/metrics-collector/pkg/adapter/clickhouse"

const (
	Version   = "1.0.0"
	APIPrefix = "/api/v1/metrics"
)

func init() {
	ApplicationConfig = &applicationConfig{}
}

type applicationConfig struct {
	MetricsAPIEnable bool              `json:"metricsApiEnable,omitempty" yaml:"metricsApiEnable"`
	Collectors       []string          `json:"collectors,omitempty" yaml:"collectors"`
	ClusterPath      string            `json:"clusterPath,omitempty" yaml:"clusterPath"`
	MetricsPath      string            `json:"metricsPath,omitempty" yaml:"metricsPath"`
	ClickHouse       clickhouse.Config `json:"clickHouse" yaml:"clickHouse"`
}

type ApiInfo struct {
	Tag         string
	Description string
}
