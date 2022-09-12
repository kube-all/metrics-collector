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

package utils

import (
	"encoding/json"
	"github.com/magiconair/properties"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog/v2"
	"strings"
)

func LoadConfig(path string, object interface{}) {
	if strings.HasSuffix(path, ".properties") {
		if props, err := properties.LoadFile(path, properties.UTF8); err != nil {
			klog.Fatalf("Unable to read application %s properties config from file, err: %s", path, err)
		} else {
			err = props.Decode(object)
			if err != nil {
				klog.Fatalf("Unable to decode application %s properties config from file, err: %s", path, err)
			}
		}
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		if data, err := ioutil.ReadFile(path); err == nil {
			if err = yaml.Unmarshal(data, object); err != nil {
				klog.Fatalf("Unable to decode application %s yaml config from file, err: %s", path, err)
			}
		} else {
			klog.Fatalf("Unable to read application %s yaml config from file, err: %s", path, err)

		}
	} else if strings.HasSuffix(path, ".json") {
		if data, err := ioutil.ReadFile(path); err == nil {
			if err = json.Unmarshal(data, object); err != nil {
				klog.Fatalf("Unable to decode application %s json config from file, err: %s", path, err)
			}
		} else {
			klog.Fatalf("Unable to read application %s json config from file, err: %s", path, err)

		}
	} else {
		klog.Fatalf("Unable to read application  config from file: %s", path)
	}
}
func StringKeyInArray(key string, arrays []string) (exist bool) {
	for _, v := range arrays {
		if key == v {
			exist = true
			return
		}
	}
	return
}