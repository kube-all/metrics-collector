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

import (
	"crypto/tls"
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2"
	"k8s.io/klog/v2"
	"time"
)

func NewClient(config Config) *sql.DB {

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: config.Address,
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.Username,
			Password: config.Password,
		},
		TLS: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			clickhouse.CompressionLZ4,
		},
		Debug: config.Debug,
	})

	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Hour)
	if err := conn.Ping(); err != nil {
		klog.Fatalf("create clickhouse connection failed, err: %s", err.Error())
	}
	return conn
}
