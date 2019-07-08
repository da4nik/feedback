FROM golang:1.11-alpine AS builder

RUN apk add bash ca-certificates git gcc g++ libc-dev
RUN go get -u github.com/da4nik/feedback/...

WORKDIR /go/src/github.com/da4nik/feedback

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/feedback

FROM alpine:3.10

RUN apk add ca-certificates

COPY --from=builder /go/bin/feedback /bin/feedback

EXPOSE 9000

ENTRYPOINT ["/bin/feedback"]
