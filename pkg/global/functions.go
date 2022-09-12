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

import (
	"errors"
	"runtime"
	"strings"
)

func RegisterApiInfo(info ApiInfo) {
	ApiInfos = append(ApiInfos, info)
}
func (a ApiInfo) Tags() []string {
	return []string{a.Tag}
}
func (app *applicationConfig) DefaultValue() {
	if !strings.HasSuffix(app.MetricsPath, "/") {
		app.MetricsPath += "/"
	}
	if runtime.GOOS == "windows" {
		app.MetricsPath = strings.Replace(app.MetricsPath, "/", "\\", -1)
	}
}
func (app *applicationConfig) Validator() (errs []error) {
	if len(app.Collectors) == 0 {
		errs = append(errs, errors.New("config's Collectors(field collectors in yaml file) must not empty"))
	}
	if len(app.ClusterPath) == 0 {
		errs = append(errs, errors.New("config's ClusterPath(field clusterPath in yaml file) must not empty"))
	}
	if len(app.MetricsPath) == 0 {
		errs = append(errs, errors.New("config's MetricsPath(field metricsPath in yaml file) must not empty"))
	}
	return errs
}
