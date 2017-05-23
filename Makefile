NETWORK=i_ipp
POSTGRES_USER=mujtaba
POSTGRES_PASSWORD=admin
POSTGRES_DB=ipp


all: run

run: main.go handlers.go model.go
	go get
	go run handlers.go main.go model.go

test:
	docker network create -d bridge --subnet 172.23.0.0/16 ipp_test
	docker run --name ipp_test_db -d \
		--network ipp_test \
		-v $(PWD)/migrations/:/docker-entrypoint-initdb.d/ \
		-e POSTGRES_DB=ipp \
		-e POSTGRES_USER=test \
		-e POSTGRES_PASSWORD=test \
		postgres:9.6
	# Wait for it to finish initing
	sleep 5
	docker run -it --rm \
		--workdir /go/github.com/mujz/ipp \
		--name ipp_test_api \
		--network ipp_test \
		-v $(PWD):/go/github.com/mujz/ipp \
		-e GOBIN=/go/bin \
		-e DB_NAME=ipp \
		-e DB_PASSWORD=test \
		-e DB_USER=test \
		-e DB_HOST=ipp_test_db \
		golang:1.8 make test-local
	docker rm -f ipp_test_db
	docker network rm ipp_test

test-local:
	go get
	go test -race -coverprofile=coverage.txt -covermode=atomic *.go

build:
	# Build OSX binary
	go build -o bin/ipp-osx
	# Build linux binary
	docker run --name ipp_build --rm -it -v $(GOPATH):/go -w /go/src/github.com/mujz/ipp golang go build -o bin/ipp *.go

