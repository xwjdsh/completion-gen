FROM golang:1.17.5-alpine3.14
Add . /go/src/completion-gen
WORKDIR /go/src/completion-gen

RUN go build ./cmd/completion-gen
