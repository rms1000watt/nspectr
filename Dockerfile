FROM golang:1.10.3-alpine3.8 AS builder
ADD . /go/src/github.com/rms1000watt/nspectr
WORKDIR /go/src/github.com/rms1000watt/nspectr
RUN apk add -U ca-certificates git && \
    go get -u github.com/kardianos/govendor && \
    govendor sync && \
    go test ./... && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags '-w -extldflags=-static' -o /nspectr

FROM scratch
COPY --from=builder /nspectr /nspectr
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/nspectr"]
