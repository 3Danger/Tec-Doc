FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder ./app .
ENV TEC_DOC_LOG_LEVEL debug
ENV TEC_DOC_FRONTEND_PORT 8002

EXPOSE 8002:8002
CMD ["./app"]
