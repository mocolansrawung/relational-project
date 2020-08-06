package main

import "net/http"

func main() {
	if err := http.ListenAndServe(":3000", ServeHTTP()); err != nil {
		panic(err)
	}
}
