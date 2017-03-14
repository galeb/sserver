build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sserver

build-mac:
	go build -o sserver

clean:
	rm -f sserver
