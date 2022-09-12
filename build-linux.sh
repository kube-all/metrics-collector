#!/usr/bin/env bash

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -o ./

docker build . -t hexiaoyun128/kube-cloud-metrics-collector:$1
docker push  hexiaoyun128/kube-cloud-metrics-collector:$1
rm ./kube-cloud-metrics-collector


