



all: notectl noteserver


.PHONY: notectl
notectl:
	go build -o ./bin/notectl ./cmd/notectl

.PHONY: noteserver
noteserver:
	go build -o ./bin/noteserver ./cmd/noteserver