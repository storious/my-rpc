package main

import (
	"fmt"
	"myRPC/example/common/model"
	"myRPC/example/consumer/proxy"
)

func main() {
	// static proxy
	//var userService service.UserService = proxy.UserServiceProxy{}
	//user := model.User{
	//	Name: "test",
	//}
	//
	//newUser := userService.GetUser(user)
	//fmt.Println(newUser)

	// dynamic proxy
	userService := proxy.NewProxy(proxy.UserService{}, proxy.ServiceProxy{})
	user := model.User{
		Name: "test",
	}
	newUser, err := userService.InvokeMethod("GetUser", user)
	if err != nil {
		panic(err)
	}
	fmt.Println(newUser.(model.User))
}
