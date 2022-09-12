


windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build  -o metrics-collector.exe  ./cmd/start.go
.PHONY: windows

macos:
	CGO_ENABLED=0    go build  -o metrics-collector ./cmd/start.go
.PHONY: macos


test:
	go test ./pkg...
.PHONY: test