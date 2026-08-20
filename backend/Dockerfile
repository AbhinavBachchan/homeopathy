FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GOTOOLCHAIN=auto
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/api

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
