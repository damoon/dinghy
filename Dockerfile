FROM golang:1.13.3-alpine3.10 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/webdav-to-s3/

FROM alpine:3.10.3
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /app/main .
CMD ["/main"]
