package main

import (
	"github.com/evermos/boilerplate-go/delivery/http"
)

func main() {
	http := http.Http{DB: db, Config: config, Router: Routes()}
	http.Serve()
}
