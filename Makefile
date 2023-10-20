all: docs

docs:
	$(MAKE) -C docs

test:
	go test -v ./...

tests: test

.PHONY: docs
