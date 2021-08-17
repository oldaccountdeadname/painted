PREFIX ?= /usr/local

all:
	go build -o build/painted *.go
install: all
	cp build/painted $(PREFIX)/bin/painted
