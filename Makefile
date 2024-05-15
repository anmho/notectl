



all: notectl noteserver gen


.PHONY: notectl
notectl:
	go build -o ./bin/notectl ./cmd/notectl

.PHONY: noteserver
noteserver:
	go build -o ./bin/noteserver ./cmd/noteserver

.PHONY: gen
gen:
	protoc --go_out=./gen --go_opt=paths=source_relative \
        --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
        ./proto/notes/note.proto
