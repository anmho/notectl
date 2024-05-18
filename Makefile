

all: notectl noteservice gen


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

.PHONY: image
image:
	./build.sh

start: noteservice
	dotenvx run -f .env.development -- ./bin/noteservice

