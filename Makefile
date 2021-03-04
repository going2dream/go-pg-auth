.PHONY: build
build: build-keygen build-api-http-server

.PHONY: build-api-http-server
build-api-http-server:
	rm -rf ./bin/api-server
	go build -v -o ./bin ./cmd/api-server

.PHONY: build-keygen
build-keygen:
	rm -rf ./bin/keygen
	go build -v -o ./bin ./cmd/keygen

check-new-keys:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: new-keys
new-keys: check-new-keys build-keygen
	./bin/keygen --output-dir=./keys