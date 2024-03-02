package service

import (
	"cloud/server"
	"cloud/server/dao"
	"cloud/server/dto"
	"cloud/server/model"
	"cloud/tool"
)

// FileService File service
type FileService struct {
	dao dao.FileDao
}

// UploadFile 文件上传
func (s *FileService) UploadFile(req dto.UploadFileDTO) (string, int) {
	// 文件信息存入数据库
	userInfo, err := server.GetUserInfoFromToken(req.Token)
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", server.SQLErrCode
	}
	id := server.GenerateUUID()
	info := model.File{
		ID:       &id,
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
	addr, err := server.UploadFile(*req.File)
	if err != nil {
		tool.Logger.Error(err.Error())
		return "", server.SQLErrCode
	}

	return addr, server.OkCode
}

// DownloadFile 文件下载
func (s *FileService) DownloadFile(req dto.DownloadFileDTO) ([]byte, int) {
	// 文件信息存入数据库
	userInfo, err := server.GetUserInfoFromToken(req.Token)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}

	fileInfo, err := s.dao.GetByID(server.GetEngine(), req.ID)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}

	// 判断是否是自己的文件
	if fileInfo.Uid != userInfo.Id {
		tool.Logger.Error(server.GetMsgByCode("zh", server.ParamErrCode))
		return nil, server.ParamErrCode
	}

	// 取文件
	file, err := server.DownloadFile(req.ID)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}

	return file, server.OkCode
}
