FROM golang:1.19 as builder

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o ingress-manager main.go   # 静态链接 否在在容器中会因为找不到动态链接库而出现 not found

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/ingress-manager .

CMD ["./ingress-manager"]
