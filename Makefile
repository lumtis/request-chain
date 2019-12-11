include Makefile.ledger
all: lint install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/rcd
		go install $(BUILD_FLAGS) ./cmd/rccli
		go install $(BUILD_FLAGS) ./cmd/rcrelayer

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify
