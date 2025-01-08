FROM golang:1.23.3-alpine as builder

WORKDIR /app

RUN apk add --no-cache build-base musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auther cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/auther /app/auther

COPY configs/config.yaml /app/configs/config.yaml

EXPOSE 5000

CMD ["/app/auther", "-config", "/app/configs/config.yaml"]
