FROM golang:1.8-alpine

ENV appDir=/go/src/svc-aggregator

RUN mkdir -p $appDir
WORKDIR $appDir

ADD . $appDir

RUN go build .

ENTRYPOINT ["/go/src/svc-aggregator/svc-aggregator"]

EXPOSE 10000
