
fmt-fetch:
	go get github.com/daixiang0/gci
	go get mvdan.cc/gofumpt

fmt: fmt-fetch
	go fmt ./...

build:
	go build -o sqlboiler-mysql drivers/sqlboiler-mysql/main.go
	go build -o sqlboiler main.go