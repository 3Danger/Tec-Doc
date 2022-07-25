FROM golang:latest as builder

RUN mkdir /build
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o app cmd/worker/main.go

FROM alpine:latest

COPY --from=builder /build/app .
CMD ["./app"]