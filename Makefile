.PHONY: build
build: build-keygen

.PHONY: build-keygen
build-keygen:
	rm -rf ./bin/keygen
	go build -v -o ./bin ./cmd/keygen

check-new-keys:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: new-keys
new-keys: check-new-keys build-keygen
	./bin/keygen -b 2048 --output-dir=./keys