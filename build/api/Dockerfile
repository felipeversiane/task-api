FROM golang:1.22-alpine3.20 as builder

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api/main.go

FROM scratch

WORKDIR /

COPY --from=builder /app/api ./

ENTRYPOINT ["./api"]
