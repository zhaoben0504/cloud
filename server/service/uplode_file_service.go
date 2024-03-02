package service

import (
	"cloud/server"
	"cloud/server/dao"
	"cloud/server/dto"
	"cloud/server/model"
	"cloud/tool"
)

// UploadFileService uploadFile service
type UploadFileService struct {
	dao dao.UploadFileDao
}

// UploadFile 文件上传
func (s *UploadFileService) UploadFile(req dto.UploadFileDTO) (string, int) {
	// 文件信息存入数据库
	userInfo, err := server.GetUserInfoFromToken(req.Token)
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", server.SQLErrCode
	}
	info := model.UploadFile{
		FileName: req.FileName,
		Path:     req.Path,
		Uid:      userInfo.Id,
	}
	err = s.dao.Add(server.GetEngine(), &info)
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", server.SQLErrCode
	}

	// 存文件
	addr, err := server.UploadCos(*req.File)
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", server.SQLErrCode
	}

	return addr, server.OkCode
}
