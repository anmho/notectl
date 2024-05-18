

all: notectl noteservice gen

.PHONY: clean
clean:
	rm ./bin/*

.PHONY: notectl
notectl:
	go build -o ./bin/notectl ./cmd/notectl

.PHONY: noteservice

noteservice:
	go build -o ./bin/noteservice ./cmd/noteservice

.PHONY: gen
gen:
	protoc --go_out=./gen --go_opt=paths=source_relative \
        --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
        ./proto/notes/note.proto

# Making image for local testing
.PHONY: image
image:
	docker build -t noteservice .

.PHONY: publish
publish:
	./build.sh

dev: noteservice
	dotenvx run -f .env.development -- ./bin/noteservice

.PHONY: env
env:
	aws secretsmanager get-secret-value \
      --secret-id noteservice-env \
      --query SecretString \
      --output text | tee .env

.PHONY: start
start:
	docker run -d --env-file .env -p 50051:50051 docker.io/anmho/noteservice:latest