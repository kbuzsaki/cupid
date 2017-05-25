.PHONY: all rall fmt test mocks

all:
	go install ./... github.com/kbuzsaki/cupid/...

rall:
	go build -a ./... github.com/kbuzsaki/cupid/...

fmt:
	gofmt -s -w -l .
	goimports -w .
 
test:
	go test -race ./...

mocks:
	rm -rf mocks
	cd server && mockery -name=Server && mv mocks ..
