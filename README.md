# boilerplate-go
## Changelog
- **v1**: first commit structure

## Description
This is an boilerplate Architecture in Go (Golang) projects.

## How To Run 

#### Run the Applications
```shell
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/evermos/boilerplate-go.git

#move to project
$ cd boilerplate-go

$ make run
```

## Live Reload
To develop the app using live reload simply execute

```shell
make dev
```

## Debugging
To debug your application on VS Code add some breakpoint particular line, then simplly hit `f5`
![Debugging Demo](assets/debugging-demo.png)


## Documentation
Swagger are automatically generated and accessible after `make run` or `make dev` on `<your-url>/swagger/index.html`

## Generating Mock ##
Mock source automatically generated when executing `make run` or `make dev`. When you add a go file and it require mock add `//go:generate go run github.com/golang/mock/mockgen -source <your-file-name>.go -destination mock/<your-file-name>_mock.go -package <your-package-name>_mock` comment right after your package declaration

example :

```go
package example

//go:generate go run github.com/golang/mock/mockgen -source service.go -destination mock/service_mock.go -package example_mock

import "github.com/gofrs/uuid"

type SomeService interface {
	ResolveByID(id uuid.UUID) (SomeEntity, error)
}

type SomeServiceImpl struct {
	SomeRepository SomeRepository `inject:"example.someRepository"`
}

func (s *SomeServiceImpl) ResolveByID(id uuid.UUID) (SomeEntity, error) {
	return s.SomeRepository.ResolveByID(id)
}
```

## Go Directories

### `configs`
The config package contains a handful of useful functions to load to configuration structs from JSON files or environment variables.

There are also many structs for common configuration options and credentials of different Cloud Services and Databases.

### `container`
Container provides an easy way to use Dependecy Injection. A container wraps those all struct into registry, if your application makes a lot of struct, you dont need to instance one by one to instantiate each struct.

### `docs`
Design and any documents including : api design,

### `internal`
This contains the business logic and the flow of data for this project only. Meaning it shouldn’t be exported outside of the project.
1. `internal/services`
    a package contains of use case layer. This main purpose to handle any business process logic.
2. `internal/repositories`
    a package contains of data storage logic handler. focus on store any data storage handler. querying, or creating/ inserting into any data storage(cache, database, any data source). This layer will act for CRUD to data storage only. No business process happen here. 
3. `internal/models`
    a package contains of entity data model that use for data storage.
4. `internal/handlers`
    a package contains of delivery layer. this main purpose to handle any I/O through the application. because main purpose this boilerplate for web api, mostly use for handle request and response to the client.
5. `internal/dto`
    a package that contains of data transform object that using for multipurpose if data need to transform to specific struct or type.

### `infras`
Any technology the project uses is written in this package, but the function written here must be a collection of functions that can be used multiple times inside multiple packages.

### `events`
The pubsub package contains two (producer and consumer) generic interfaces for producing data to queues as well as subscribing and consuming data from those queues.

### `transport`
Transport contains helpers applicable to all supported transports, `http, grpc, or etc`.
1. `transport/http`
    Package http provides a general purpose HTTP binding for endpoints.

### `shared`
a package that helps for the project, whether it is response, constant variable or other.

### `router`
router is a package that is used to wrap all handlers together, then use it on the required transport method.

### `bootstrap`
System init (systemd, upstart, sysv) and process manager/supervisor (runit, supervisord) configs

## Tools References
- [Environment Configuration](https://github.com/spf13/viper)
- [Router](https://github.com/go-chi/chi)
- [Dependency Injection](https://github.com/facebookarchive/inject)
- [Live Reload](https://github.com/cosmtrek/air)
- [Mocks](https://github.com/golang/mock)
- [Event Consumer](https://github.com/nsqio/go-nsq)
- [Linter - golangci-lint](github.com/golangci/golangci-lint)
- [Documentation Generator - Swaggo](github.com/swaggo/swag/cmd/swag)
- [Documentation Server - Swaggo HTTP Server](github.com/swaggo/http-swagger)

## Useful Links
- [Golang Tools Best Practice](https://github.com/go-modules-by-example/index/tree/master/010_tools)

## Diagram
Here's the diagram to explanation about project structre.
[Diagram](https://drive.google.com/file/d/1uxb2dwHA1GFWuPs9ljhsMBMCjP6gVGtk/view)