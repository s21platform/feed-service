FROM golang:1.24 as builder

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/main cmd/service/main.go
RUN go build -o build/worker_user cmd/workers/user/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /usr/src/service/build/main /app
COPY --from=builder /usr/src/service/build/worker_user /app
RUN apk add --no-cache gcompat
RUN chmod +x main

CMD ./main & ./worker_user