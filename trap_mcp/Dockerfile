FROM golang:1.24-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ --mount=type=cache,target=/root/.cache/go-build/ \
    CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o /app/main .

FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD [ "./main" ]
