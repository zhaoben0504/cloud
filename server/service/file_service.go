package service

import (
	"cloud/server"
	"cloud/server/dao"
	"cloud/server/dto"
	"cloud/server/model"
	"cloud/server/vo"
	"cloud/tool"
	"errors"
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

// DeleteFile 文件删除
func (s *FileService) DeleteFile(req dto.DeleteFileDTO) int {
	file, err := s.dao.GetByID(server.GetEngine(), req.ID)
	if err != nil {
		tool.Logger.Error(err.Error())
		return server.SQLErrCode
	}

	userInfo, err := server.GetUserInfoFromToken(req.Token)
	if err != nil {
		tool.Logger.Error(err.Error())
		return server.InternalErrCode
	}

	// 判断删除的文件是不是自己的
	if *userInfo.Id != *file.Uid {
		tool.Logger.Error(errors.New(server.GetMsgByCode("zh", server.PermissionErrCode)))
		return server.PermissionErrCode
	}

	// 删除数据库中文件信息
	timeStamp := tool.UnixMillisecond()
	err = s.dao.EditByID(server.GetEngine(), req.ID, &model.File{DeletedAt: &timeStamp})
	if err != nil {
		tool.Logger.Error(err.Error())
		return server.InternalErrCode
	}

	err = server.DeleteFile(req.ID)
	if err != nil {
		tool.Logger.Error(err.Error())
		return server.InternalErrCode
	}

	return server.OkCode
}

// ListFile 文件查询
func (s *FileService) ListFile(req dto.ListFileDTO) (*vo.FileListVO, int) {
	where := "deleted=0"
	param := make([]interface{}, 0)
	if req.Keyword != "" {
		where += " AND `file_name` like ?"
		param = append(param, "%s"+req.Keyword+"%s")
	}

	sort := make(map[string][]string)
	sort["desc"] = []string{"created_at"}
	list, total, err := s.dao.GetList(server.GetEngine(), where, param, req.Page, req.Rows, sort)
	if err != nil {
		tool.Logger.Error(err.Error())
		return nil, server.SQLErrCode
	}

	return &vo.FileListVO{Total: total, List: list}, server.OkCode
}
