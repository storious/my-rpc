package service

import "myRPC/example/common/model"

type UserService interface {
	GetUser(user model.User) model.User
}
