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

package main

import (
	"flag"
	"os"

	"github.com/kube-all/metrics-collector/cmd/collector"
	"github.com/kube-all/metrics-collector/cmd/metrics"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

func commandRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "metrics-collector",
		Long: `metrics collector`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(2)
		},
	}
	rootCmd.Flags().SortFlags = true
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("logtostderr"))
	pflag.CommandLine.Set("logtostderr", "true")

	// add sub command
	rootCmd.AddCommand(collector.KubeMetricsCollectorCommand())
	rootCmd.AddCommand(metrics.KubeMetricsGenerateCommand())
	return rootCmd
}

func main() {
	defer klog.Flush() // flushes all pending log I/O
	command := commandRoot()
	command.Execute()
}
