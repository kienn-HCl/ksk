BINARY := ksk

.PHONY: build clean snapshot check

build:
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY) result
	rm -rf dist/

snapshot:
	goreleaser release --snapshot --clean

check:
	goreleaser check
