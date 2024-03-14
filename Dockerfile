# syntax=docker/dockerfile:1
FROM golang:1.22.1 as builder

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main ./cmd/app/main.go
RUN go test -race ./...

FROM alpine as runner

COPY --from=builder /app/main /app/main

EXPOSE 8000
CMD ["/app/main"]
