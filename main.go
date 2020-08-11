package main

import (
	"github.com/evermos/boilerplate-go/delivery/http"
)

func main() {
	registry()
	http := http.Http{DB: db, Config: config, Router: router.NewRouter()}
	http.Serve()
}
