BINARY=engine
test: clean documents generate
	go test -v -cover -covermode=atomic ./...

coverage: clean documents generate
	bash coverage.sh --html
	
engine:
	go build -o ${BINARY} *.go

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@find . -name *mock* -delete
	@rm -rf .cover

docker:
	docker build -t boilerplate-go -f Dockerfile-local .

run:
	docker-compose up --build -d

stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

lint:
	./bin/golangci-lint run ./...

documents:
	swag init

generate:
	go generate ./...
	
.PHONY: test coverage engine clean build docker run stop lint-prepare lint documents generate
