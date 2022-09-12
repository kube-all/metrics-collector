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

package clickhouse

type Config struct {
	Address            []string `json:"address,omitempty" yaml:"address"`
	Database           string   `json:"database,omitempty" yaml:"database"`
	Username           string   `json:"username,omitempty" yaml:"username"`
	Password           string   `json:"password,omitempty" yaml:"password"`
	InsecureSkipVerify bool     `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify"`
	Debug              bool     `json:"debug" yaml:"debug"`
}
