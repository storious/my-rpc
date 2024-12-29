package main

import (
	"myRPC/example/provider"
	"myRPC/example/provider/server"
	"myRPC/rpc-easy/registry"
)

func main() {
	httpServer := server.NewHttpServer()
	var r = registry.NewRegistry(registry.GetLocalCache())
	r.Register("UserService", provider.UserLogic{})
	err := httpServer.Start(":8080")
	if err != nil {
		return
	}
}
