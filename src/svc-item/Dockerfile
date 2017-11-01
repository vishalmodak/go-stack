FROM golang:1.8-alpine

ENV appDir=/go/src/svc-item

RUN mkdir -p $appDir
WORKDIR $appDir

ADD . $appDir

RUN go build .

CMD ["/go/src/svc-item/svc-item"]

EXPOSE 10001
