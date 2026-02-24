.PHONY: run build clean test vet

BINARY := poker_tui
CMD     := ./cmd/poker_tui

run:
	go run $(CMD)

build:
	go build -o $(BINARY) $(CMD)

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64 $(CMD)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY).exe $(CMD)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY) $(BINARY)-linux-amd64 $(BINARY).exe
