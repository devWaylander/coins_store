FROM golang:1.22.8 as builder

WORKDIR /usr/src/app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 go build -o /usr/src/app/coins ./cmd

FROM alpine:latest

COPY --from=builder /usr/src/app/coins /coins

EXPOSE 8080

CMD ["/coins"]