FROM golang:1.21 as builder

WORKDIR /sasposts
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /sasposts
COPY --from=builder /sasposts/server .
COPY --from=builder /sasposts/migrations ./migrations

EXPOSE 8080
CMD ["./server"]