```shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct
go get k8s.io/api/core/v1@v0.23.3
go get k8s.io/client-go@v0.24.3
```
