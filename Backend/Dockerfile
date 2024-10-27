FROM golang:alpine AS builder

WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY .env .
RUN go build -o /go/bin/app ./cmd/http

FROM alpine 
COPY --from=builder /go/bin/app /app
COPY --from=builder /go/src/.env /
ENTRYPOINT ["/app"]
