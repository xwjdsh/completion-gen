FROM golang:1.17.5-alpine3.15 AS builder
Add . /go/src/completion-gen
WORKDIR /go/src/completion-gen
RUN go build ./cmd/completion-gen

FROM alpine:3.15.0
WORKDIR /workspace
COPY --from=builder /go/src/completion-gen/completion-gen /workspace
COPY --from=builder /go/src/completion-gen/tmpls /workspace/tmpls
COPY --from=docker:20.10.12-dind-alpine3.15 /usr/local/bin/docker /usr/local/bin
ENTRYPOINT ["/workspace/completion-gen"]

