FROM arm64v8/golang:1.16 AS builder

RUN apt-get update && apt-get install -y gcc-aarch64-linux-gnu

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build 

RUN go test

FROM arm64v8/golang:1.16 AS runner

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]