package main

import (
	"myRPC/example/provider"
	"myRPC/rpc-easy/registry"
	"myRPC/rpc-easy/server"
)

// TODO: provide user service example
func main() {
	httpServer := server.NewHttpServer()
	registry.Register("UserService", &provider.UserLogic{})
	err := httpServer.Start(":8080")
	if err != nil {
		return
	}
}
