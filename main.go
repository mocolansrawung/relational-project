package main

func main() {
	if err := ServeHTTP(); err != nil {
		panic(err)
	}
}
