deps:
	@if [ -z `which dep` ]; then \
    go get -u github.com/golang/dep/cmd/dep; \
	fi; \
	dep ensure -vendor-only

lint:
	go vet ./...

test:
	go test -race -cover -v ./... -timeout 90

agent:
	go build -o build/agent github.com/nicktitle/redcanary/cmd
