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
package apis

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	v1 "github.com/kube-all/metrics-collector/pkg/apis/v1"
	"github.com/kube-all/metrics-collector/pkg/embeds"
	"github.com/kube-all/metrics-collector/pkg/global"
	"k8s.io/klog/v2"
)

func AddResources() {
	klog.Info("add service resources")
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	container := restful.DefaultContainer
	container.Filter(container.OPTIONSFilter)
	v1.MetricsResource{}.AddRestApiService(container)
	klog.Info("add resource success")
}
func AddSwagger() {
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// static file serving
	//http.Handle("/swagger/", http.FileServer(swagger.StaticFileSystem()))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(embeds.StaticFileSystem())))
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Metrics Collector",
			Description: "Kubernetes Prometheus Metrics Collector",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "kubeall",
					Email: "kubeall@alimai.com",
					URL:   "",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{}
	for _, s := range global.ApiInfos {
		props := spec.TagProps{Name: s.Tag, Description: s.Description}
		swo.Tags = append(swo.Tags, spec.Tag{TagProps: props})
	}
}
