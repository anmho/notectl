FROM golang:1.22 AS builder
WORKDIR /app

COPY . .

RUN apt update
RUN apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN export PATH="$PATH:$(go env GOPATH)/bin"

RUN make gen

RUN go mod tidy
RUN make noteservice

CMD ["./bin/noteservice"]



# Image step
#FROM alpine:3.19.1
#
#WORKDIR /app
#COPY --from=builder /app/bin/noteservice /app/noteservice
#
#CMD ["./noteservice"]