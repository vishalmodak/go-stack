FROM golang:1.8-alpine

ENV appDir=/go/src/svc-price

RUN mkdir -p $appDir
WORKDIR $appDir

ADD . $appDir

RUN go build .

CMD ["/go/src/svc-price/svc-price"]

EXPOSE 10002
