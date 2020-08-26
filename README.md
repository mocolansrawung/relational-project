# boilerplate-go
## Changelog
- **v1**: first commit structure

## Description
This is an boilerplate Architecture in Go (Golang) projects.

## How To Run 

#### Run the Applications
```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/evermos/boilerplate-go.git

#move to project
$ cd boilerplate-go

# change .env.example to .env
$ mv .env.example .env

# Build the docker image
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker ps

# Execute the call
$ curl localhost:9090

# Stop
$ make stop
```

## Ez Development
#### Refresh

Use refresh are for better development experience. 

#
```bash
#install on your GO ROOT/ GO PATH
$ go get github.com/markbates/refresh

#Create config file
$ refresh init

#Run like a boss
$ refresh run
```

## Documentation ##
#### Swagger
#
```bash
$ go get
$ swag init
```

go to : {app_url}/docs/

example:
```
localhost:8000/docs/
```

## Generating Mock ##
```
# generate for service
$ mockgen -source=internals/services/service.go -package=mocks -destination=mocks/service.go -mock_names="ExampleContract=MockExampleServiceContract"

# generate for repository
$ mockgen -source=internals/repositories/repositories.go -package=mocks -destination=mocks/repositories.go -mock_names="ExampleContract=MockExampleRepositoryContract"
```
