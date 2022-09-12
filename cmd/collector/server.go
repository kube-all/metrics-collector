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
package collector

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/kube-all/metrics-collector/cmd/collector/options"
	"github.com/kube-all/metrics-collector/pkg/apis"
	"github.com/kube-all/metrics-collector/pkg/collectors"
	"github.com/kube-all/metrics-collector/pkg/embeds"
	"github.com/kube-all/metrics-collector/pkg/global"
	"github.com/kube-all/metrics-collector/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
)

func KubeMetricsCollectorCommand() *cobra.Command {
	s := options.NewKubeMetricsCollectorOptions()
	cmd := &cobra.Command{
		Use:     "collector",
		Short:   "start metrics collector server",
		Version: global.Version,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Run(s); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	s.AddFlag(cmd.Flags())
	return cmd
}

func Run(s *options.KubeMetricsCollectorOptions) (err error) {
	utils.LoadConfig(s.Config, global.ApplicationConfig)
	global.ApplicationConfig.DefaultValue()
	if errs := global.ApplicationConfig.Validator(); len(errs) > 0 {
		for _, er := range errs {
			klog.Error(er.Error())
		}
		os.Exit(-1)
	}
	// generate default metrics config
	generate(global.ApplicationConfig.MetricsPath)
	// load metrics config
	storage := collectors.GetCollectorStorageInstance(global.ApplicationConfig.ClusterPath,
		global.ApplicationConfig.MetricsPath, global.ApplicationConfig.Collectors)

	go storage.Watcher(global.ApplicationConfig.Collectors)
	if cfg, err := clientcmd.BuildConfigFromFlags(s.MasterUrl, s.KubConfig); err == nil {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		kubeClient := kubernetes.NewForConfigOrDie(cfg)
		rl, _ := resourcelock.New(
			resourcelock.ConfigMapsResourceLock,
			"",
			"metrics-collector-lock",
			kubeClient.CoreV1(),
			kubeClient.CoordinationV1(),
			resourcelock.ResourceLockConfig{
				Identity: "metrics-collector-lock",
			},
		)
		go leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
			Lock:            rl,
			ReleaseOnCancel: true,
			LeaseDuration:   15 * time.Second,
			RenewDeadline:   10 * time.Second,
			RetryPeriod:     2 * time.Second,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {

				},
				OnStoppedLeading: func() {
					klog.Warning("leader change")
				},
				OnNewLeader: func(identity string) {
				},
			},
		})

	} else {

	}
	apis.AddResources()
	apis.AddSwagger()
	klog.Info("read to start http server")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		klog.Errorf("start http server failed, err: %s", err.Error())
	}
	klog.Info("application shutting down")
	return err
}
func generate(dataPath string) {
	var (
		files []string
	)
	filepath.Walk(dataPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			klog.Warningf("load data from path: %s failed, err: %s", path, err.Error())
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(dataPath, path)
		if err != nil {
			klog.Warningf("file path: %s rel: %s, err: %s ", path, dataPath, err.Error())
			return nil
		}

		if strings.HasSuffix(relPath, ".yaml") || strings.HasSuffix(relPath, ".yml") {
			files = append(files, relPath)
		}
		return nil
	})
	if len(files) > 0 {
		for _, f := range files {
			klog.Infof("will use metrics file: %s ", f)
		}
	} else {
		klog.Infof("will use default metrics file and will write to %s", dataPath)
		_, err := os.Stat(dataPath)
		if os.IsNotExist(err) {
			os.MkdirAll(dataPath, os.ModePerm)
		}
		files, _ := embeds.MetricsYamlFiles.ReadDir("metrics")
		for _, f := range files {
			if data, e := embeds.MetricsYamlFiles.ReadFile(path.Join("metrics", f.Name())); e == nil {
				fp := path.Join(dataPath, f.Name())
				if e1 := ioutil.WriteFile(fp, data, os.ModePerm); e1 == nil {
					klog.Infof("generate file: %s success", fp)
				} else {
					klog.Errorf("generate file: %s failed, err: %s", fp, e1.Error())
				}
			}
		}
	}
}
