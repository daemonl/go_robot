

lint:
	go vet ./...
	golint ./...

test: lint
	go test -v -cover

run:
	go run ./cmd/robot/*.go --quiet --max-x 4 --max-y 4

runspec/%:
	cat ./example/$*.txt | go run ./cmd/robot/*.go --quiet --max-x 4 --max-y 4

build:
	go build -o ./bin/robot ./cmd/robot/*.go


.PHONY: lint test run build


