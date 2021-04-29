FROM golang:1.16.1

RUN mkdir -p /go/src/github.com/ksivvi0/statusdetector

WORKDIR /go/src/github.com/ksivvi0/statusdetector

ADD . /go/src/github.com/ksivvi0/statusdetector

RUN go install github.com/ksivvi0/statusdetector

EXPOSE 8001