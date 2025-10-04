FROM golang:1.24.5-bullseye AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -trimpath \
    -ldflags="-s -w -extldflags=-static" \
    -o /out/app .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/app /app
ENTRYPOINT ["/app/app"]
