.PHONY: all rall fmt test mocks fresh

all:
	go install github.com/kbuzsaki/cupid/cmd/...

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

fresh:
	rm -rf raftexample-*

launch1: fresh all
	./launch.sh 1

launch3: fresh all
	./launch.sh 3

launch5: fresh all
	./launch.sh 5
