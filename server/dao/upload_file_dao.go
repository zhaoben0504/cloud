package dao

import (
	"cloud/server"
	"cloud/server/model"
	"cloud/tool"
	"errors"
	"github.com/go-xorm/xorm"
)

type UploadFileDao struct {
}

// Add 插入单条记录
func (d UploadFileDao) Add(engine *xorm.EngineGroup, entity *model.UploadFile) error {
	if engine == nil {
		return errors.New(server.GetMsgByCode("en", server.InternalErrCode))
	}

	time := tool.UnixSecond()
	if entity.CreatedAt == nil {
		entity.CreatedAt = &time
	}
	if entity.UpdatedAt == nil {
		entity.UpdatedAt = &time
	}
	if entity.DeletedAt == nil {
		entity.DeletedAt = new(int64)
	}
	_, err := engine.Insert(entity)
	if err != nil {
		return err
	}
	return nil
}
