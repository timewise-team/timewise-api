package main

import (
	_ "api/docs"
	"api/server"
)

// @title timewise-api
// @version 1.0
// @description Timewise API
// @in header
func main() {
	server.RegisterServer()
}
