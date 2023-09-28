FROM golang:alpine AS builder
WORKDIR /build
COPY . /build/
RUN go mod tidy && go build -o gobot cmd/main.go

FROM scratch
COPY --from=builder /build/gobot /
EXPOSE 5000
ENTRYPOINT ["/gobot"]