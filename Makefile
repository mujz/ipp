NETWORK=i_ipp
POSTGRES_USER=mujtaba
POSTGRES_PASSWORD=admin
POSTGRES_DB=ipp


all: run

run: main.go handlers.go model.go
	go get ./...
	go run handlers.go main.go model.go

build:
	# Build OSX binary
	go build -o bin/ipp-osx
	# Build linux binary
	docker run --name ipp_build --rm -it -v $(GOPATH):/go -w /go/src/github.com/mujz/i++ golang go build -o bin/ipp *.go

