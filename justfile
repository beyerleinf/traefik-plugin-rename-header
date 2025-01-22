lint:
  golangci-lint run

test:
	go test -v -cover ./...

vendor:
	go mod vendor

clean:
	rm -rf ./vendor