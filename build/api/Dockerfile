FROM golang:1.22-alpine3.20 as builder

WORKDIR /app
COPY . .

ENV GIN_MODE=${GIN_MODE:-release}

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go

FROM scratch
WORKDIR /
COPY --from=builder /app/api ./

ENTRYPOINT ["./api"]