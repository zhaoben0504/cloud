package service

import (
	"cloud/server"
	"cloud/server/dao"
	"cloud/server/dto"
	"cloud/server/model"
	"cloud/tool"
	"context"
	"fmt"
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

	// 2.生成token并塞入Redis
	*token, err = server.GenerateToken(user)
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

func (s *UserService) EmailCode(req *dto.EmailCodeDTO) (token *string, code int) {
	// 1.从数据库读取数据
	user := new(model.UserBasic)
	has, err := server.GetEngine().Where("email= ?", req.Email).Get(user)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}
	if !has {
		return nil, server.ParamErrCode
	}

	// 2.邮箱存在
	emailCode := server.GenerateEmailCode()
	err = server.SetInfoInRedis(req.Email, emailCode, server.CodeExprie)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.InternalErrCode
	}

	data := fmt.Sprintf("已经向%s的邮箱成功发送验证码, 验证码是%s", req.Email, emailCode)

	return &data, server.OkCode
}

func (s *UserService) UserInfo(req *dto.UserInfoDTO) (userInfo *model.UserBasic, code int) {
	// 1.从数据库读取数据
	//user := new(model.UserBasic)
	err := server.GetEngine().Where("id= ?", req.Id).Find(userInfo)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}
	if userInfo == nil {
		return nil, server.UserNotExistErrCode
	}

	return userInfo, server.OkCode
}
