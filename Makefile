.PHONY: all rall fmt test

all:
	go install ./... github.com/kbuzsaki/cupid/...

rall:
	go build -a ./... github.com/kbuzsaki/cupid/...

fmt:
	gofmt -s -w -l .
 
test:
	go test ./...
