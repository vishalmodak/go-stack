FROM golang:1.8-alpine

ENV appDir=/go/src/svc-item

RUN mkdir -p $appDir
WORKDIR $appDir

ADD . $appDir

RUN go build .

ENTRYPOINT ["/go/src/svc-aggregator/svc-item"]

EXPOSE 10001
