.PHONY: test

test:
	go test ./... -coverprofile=coverage.out ./...
	go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]