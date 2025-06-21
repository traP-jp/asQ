FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go build trap_mcp/main.go

FROM gcr.io/distroless/static-debian12:latest
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
