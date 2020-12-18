FROM golang:1.14 as build

ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on

WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o app main.go

FROM ubuntu:18.04 as prod
EXPOSE 5000
COPY --from=build /go/release/app /
COPY --from=build /go/release/config /config
CMD ["/app"]