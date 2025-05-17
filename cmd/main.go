package main

import "github.com/tranluan258/go-tcp/server"

func main() {
	server := server.New()
	server.Run()
}
