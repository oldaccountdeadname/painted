all:
	go build -o build/painted *.go
install: all
	cp build/painted /usr/local/bin/painted
