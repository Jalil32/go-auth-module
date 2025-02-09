ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-bookworm as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN ls -al /app
RUN go mod download && go mod verify

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/app

# Create a minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy built binary
COPY --from=builder /app/server .
# TODO: Remove this at some stage. I don't think its secure
#       Should be using fly secrets instead for production...
COPY --from=builder /app/.env .

EXPOSE 3000

# Run the binary.
CMD ["./server"]
