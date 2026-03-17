FROM golang:1.26-alpine AS builder

ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags="-s -w" -trimpath -o qqbot . && chmod +x qqbot

FROM alpine

WORKDIR /root/

COPY --from=builder /app/qqbot .
COPY --from=builder /app/configs/config.toml .

CMD ["./qqbot", "-config", "config.toml"]