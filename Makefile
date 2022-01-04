.PHONY: build verify test

all: build

build:
	go build 

test:
	go test ./...

clean:
	rm duckcoin
	rm blockchain.db
	rm duckcoin.wallet

