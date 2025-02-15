FROM golang:1.22.8 as builder

WORKDIR /usr/src/app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /usr/src/app/coins ./cmd

FROM alpine:3.21.0

RUN apk update && \
    apk add --no-cache bash curl && \
    curl -Lo /usr/local/bin/wait-for-it https://github.com/vishnubob/wait-for-it/releases/download/v2.8.0/wait-for-it-linux-amd64 && \
    chmod +x /usr/local/bin/wait-for-it

COPY --from=builder /usr/src/app/coins /coins

EXPOSE 8080

CMD ["/usr/local/bin/wait-for-it", "${DB_HOST}:${DB_PORT}", "--", "/coins"]