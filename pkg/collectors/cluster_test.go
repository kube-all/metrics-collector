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
package collectors

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestClusters(t *testing.T) {
	defs := []Cluster{
		{
			Name: "demo-jwt",
			Code: "demo-jwt",
			Prometheus: Prometheus{
				Address: "http://prometheus.com",
				Token:   "jwt token",
			},
		},
		{
			Name: "demo-basic",
			Code: "demo-basic",
			Prometheus: Prometheus{
				Address:  "http://prometheus.com",
				Username: "basic username",
				Password: "basic password",
			},
		},
	}
	data, err := yaml.Marshal(defs)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("../embeds/clusters.yaml", data, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
