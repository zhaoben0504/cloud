package service

import "cloud/server/dao"

// UserService User service
type UserService struct {
	dao dao.UserDao
}
