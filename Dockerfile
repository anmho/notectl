# Build stage
FROM --platform=linux/amd64 golang:1.22 AS builder
WORKDIR /app

COPY . .

# Install compilation dependencies
RUN apt update
RUN apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN go mod tidy

# Compile
RUN export PATH="$PATH:$(go env GOPATH)/bin"
RUN make gen
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/noteservice ./cmd/noteservice/main.go

# Final step
FROM --platform=linux/amd64 alpine:3.19.1

WORKDIR /app
# Copy from builder
COPY --from=builder /app/bin/noteservice /app/noteservice

CMD ["./noteservice"]