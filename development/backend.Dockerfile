FROM golang:1.22-alpine AS builder

RUN apk add --no-cache build-base
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /out/suuq-server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/suuq-server /app/suuq-server

EXPOSE 8080
CMD ["/app/suuq-server"]
