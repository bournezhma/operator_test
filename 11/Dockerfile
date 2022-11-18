FROM golang:1.19 as builder

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLE=0 go build -o ingress-manager main.go

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/ingress-manager .

CMD ["./ingress-namager"]