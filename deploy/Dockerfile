
FROM golang:alpine AS builder
WORKDIR /build
COPY . /build
RUN apk --no-cache add ca-certificates
RUN go mod download && go mod verify
RUN go build -o gobot cmd/gobot.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/gobot /
EXPOSE 5000
ENTRYPOINT ["/gobot"]