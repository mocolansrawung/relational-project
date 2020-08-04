package main

import (
	"fmt"

	"github.com/evermos/boilerplate-go/configs"
)

func init() {
	config = configs.Get()
	initDb()
}

func main() {
	fmt.Println("Hello World")
}
