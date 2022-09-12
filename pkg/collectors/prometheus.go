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
	"crypto/tls"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"net/http"
	"strings"
	"time"
)

type ClusterPrometheus struct {
	Cluster Cluster
}

func (c *ClusterPrometheus) NewClient() (api v1.API, err error) {
	var (
		tr     *http.Transport
		client api.Client
	)
	if strings.HasPrefix(c.Cluster.Prometheus.Address, "https://") {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{}
	}

	client, err = api.NewClient(api.Config{
		Address:      c.Cluster.Prometheus.Address,
		RoundTripper: tr,
	})
	if err != nil {
		return
	}
	v1api := v1.NewAPI(c.Client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
}
