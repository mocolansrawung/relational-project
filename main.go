package main

import (
	"github.com/evermos/boilerplate-go/delivery/http"
)

func main() {
	registry()
	http := http.Http{DB: db, Config: config, Router: router.Route()}
	http.Serve()
}
