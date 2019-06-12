FROM golang

LABEL maintainer="michal.jemala@gmail.com"

ENV GO111MODULE on

WORKDIR /go/src/payments-sample

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install ./...

EXPOSE 80
