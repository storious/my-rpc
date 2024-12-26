package provider

import (
	"fmt"
	"myRPC/example/common/model"
)

type UserLogic struct {
}

func (u UserLogic) GetUser(user model.User) model.User {
	fmt.Println("user name:", user.Name)
	return user
}
