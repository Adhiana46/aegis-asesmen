# Builder stage
FROM golang:1.22-alpine as builder

RUN apk --no-cache add tzdata

RUN mkdir /app

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download


# Copy files
COPY cmd ./cmd
COPY config ./config
COPY database ./database
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 go build -o api-service ./cmd/api

RUN chmod +x api-service

# Runner stage
FROM scratch

# Copy the timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=UTC

# Copy the built binary
COPY --from=builder /app/api-service /api-service

# Copy database migrations
COPY --from=builder /app/database /database

WORKDIR /

EXPOSE 9000

ENTRYPOINT ["/api-service"]