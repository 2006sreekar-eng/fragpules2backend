FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy everything from the backend folder context
COPY . .

# Generate dependencies and build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /fragpulse-server ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /fragpulse-server .

# Create a migrations folder so the app doesn't crash if it looks for it
RUN mkdir -p ./migrations

EXPOSE 8080

CMD ["./fragpulse-server"]