FROM golang:1.22.8 as builder

WORKDIR /usr/src/app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /usr/src/app/coins ./cmd

FROM alpine:3.21.0

RUN apk add --no-cache curl
RUN curl -fsSL -o /usr/local/bin/wait-for https://github.com/eficode/wait-for/releases/download/v2.2.4/wait-for
RUN chmod +x /usr/local/bin/wait-for
RUN ls -l /usr/local/bin/wait-for

COPY --from=builder /usr/src/app/coins /coins

EXPOSE 8080

CMD ["/usr/local/bin/wait-for", "${DB_HOST}:${DB_PORT}", "--", "/coins"]