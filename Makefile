sserver: clean
	type godep > /dev/null 2>&1 || go get github.com/tools/godep
	godep restore ./...
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build

clean:
	rm -f sserver

