package main

import (
	"fmt"
	"myRPC/example/common/model"
	"myRPC/example/common/service"
	"myRPC/example/consumer/proxy"
)

// TODO: consumer example
func main() {
	var userService service.UserService = proxy.UserServiceProxy{}

	user := model.User{
		Name: "test",
	}

	newUser := userService.GetUser(user)
	fmt.Println(newUser)
}
