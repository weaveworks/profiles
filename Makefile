build:
	go build -o ./pctl cmd/pctl/main.go

test: unit-test integration-test

unit-test:
	go test -count=1 ./pkg/...

integration-test:
	go test -count=1 ./tests/...

lint:
	golangci-lint run
