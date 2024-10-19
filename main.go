package main

import (
	_ "api/docs"
	"api/server"
)

// @title timewise-api
// @version 1.0
// @description Timewise API
// @securityDefinitions.apiKey bearerToken
// @in header
// @name Authorization
func main() {
	server.RegisterServer()
}
