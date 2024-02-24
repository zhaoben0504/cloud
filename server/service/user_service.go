package service

import (
	"cloud/server"
	"cloud/server/dao"
	"cloud/server/dto"
	"cloud/server/model"
	"cloud/tool"
	"context"
)

// UserService User service
type UserService struct {
	dao dao.UserDao
	ctx context.Context
}

func (s *UserService) Login(req *dto.UserLoginDTO) (token *string, code int) {
	// 1.从数据库读取数据
	user := new(model.UserBasic)
	has, err := server.GetEngine().Where("name= ? AND password= ?", req.Name, tool.Md5(req.Password)).Get(user)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}
	if !has {
		return nil, server.ParamErrCode
	}
	// 2.返回token
	*token, err = tool.GenerateToken(*user.Id, *user.Identity, *user.Name, int64(server.TokenExpire))
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.InternalErrCode
	}

	return token, server.OkCode
}

func (s *UserService) Register(req *dto.UserRegisterDTO) (*int64, int) {
	code, err := server.GetRedisClient().Get(s.ctx, req.Email).Result()
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.VerificationCodeErrCode
	}
	if code != req.Code {
		return nil, server.VerificationCodeErrCode
	}
	// 验证结束， 判定用户是否存在
	has, err := server.GetEngine().Where("name= ? and deleted_at > 0", req.Name).Get(&model.UserBasic{})
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}

	if has {
		return nil, server.UserAlreadyExistErrCode
	}

	// 写入用户
	md5Pwd := tool.MD5(req.Password)
	uc := &model.UserBasic{
		Name:     &req.Name,
		Password: &md5Pwd,
		Email:    &req.Email,
	}
	id, err := s.dao.Add(server.GetEngine(), uc)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}

	return &id, server.OkCode
}
