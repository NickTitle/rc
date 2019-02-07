deps:
	@if [ -z `which dep` ]; then \
    go get -u github.com/golang/dep/cmd/dep; \
	fi; \
	dep ensure -vendor-only

lint:
	go vet ./...

test:
	go test -race -cover -v ./... -timeout 90s

agent:
	go build -o build/agent github.com/nicktitle/rc/cmd
