FROM alpine:3.10

WORKDIR /work


COPY ./kube-cloud-metrics-collector /work

EXPOSE 8080

ENTRYPOINT ["/work/kube-cloud-metrics-collector"]