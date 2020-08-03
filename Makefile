BINARY=engine
engine:
	go build -o ${BINARY} main.go

unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t boilerplate-go -f Dockerfile .

run:
	docker-compose up --build -d

stop:
	docker-compose down

.PHONY: engine clean unittest build docker run stop