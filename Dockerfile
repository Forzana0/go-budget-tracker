FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/budget-tracker main.go

FROM alpine:3.23

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/bin/budget-tracker /app/budget-tracker

COPY --from=builder /app/migrations /app/migrations

EXPOSE 9091

CMD ["/app/budget-tracker"]
