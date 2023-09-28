FROM golang:alpine AS builder
WORKDIR /build
COPY . /build/
RUN apk --no-cache add ca-certificates
RUN go mod tidy && go build -o gobot cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/gobot /
EXPOSE 5000
ENTRYPOINT ["/gobot"]