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

package v1

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/kube-all/metrics-collector/pkg/global"
)

type MetricsResource struct {
}

func (c MetricsResource) AddRestApiService(container *restful.Container) {
	ws := new(restful.WebService)
	apiInfo := global.ApiInfo{}
	apiInfo.Tag = "metrics"
	apiInfo.Description = "集群Metrics"
	global.RegisterApiInfo(apiInfo)
	ws.Path(global.APIPrefix + "/namespaces").
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").
		To(c.list).
		Doc("集群Namespace Metrics列表").
		Param(ws.QueryParameter("page", "翻页查询第几页").DataType("number")).
		Param(ws.QueryParameter("limit", "每页数量").DataType("number")).
		Param(ws.QueryParameter("namespace", "Namespace").DataType("string")).
		Param(ws.QueryParameter("clusterId", "Cluster ID").DataType("number")).
		Metadata(restfulspec.KeyOpenAPITags, apiInfo.Tags()))

	container.Add(ws)
}

func (c MetricsResource) list(req *restful.Request, resp *restful.Response) {
}
