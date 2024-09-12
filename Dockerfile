# syntax=docker/dockerfile:1
FROM golang:1.22 AS builder

WORKDIR /app

# Install dependencies
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

# Compile binary
COPY ./ ./
RUN CGO_ENABLED=0 go build -o ./main ./cmd/app/main.go

FROM alpine:3.19 AS runner

COPY --from=builder /app/main /app/main

EXPOSE 8000
CMD ["/app/main"]
