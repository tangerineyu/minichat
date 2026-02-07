.PHONY: wire test

wire:
	go install github.com/google/wire/cmd/wire@v0.7.0
	cd internal/di && wire

test:
	go test ./...

