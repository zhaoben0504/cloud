package dao

import (
	"cloud/server"
	"cloud/server/model"
	"cloud/tool"
	"errors"
	"github.com/go-xorm/xorm"
)

type UserDao struct {
}

// Add 插入单条记录
func (d UserDao) Add(engine *xorm.EngineGroup, entity *model.UserBasic) (int64, error) {
	if nil == engine {
		tool.Logger.Error("engine is empty")
		return -1, errors.New("engine is empty")
	}

	if entity.Id == nil {
		id := server.GetID()
		entity.Id = &id
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
		tool.Logger.Error(err)
		return -1, err
	}
	return *entity.Id, nil
}

// GetByID 根据id查询
func (d UserDao) GetByID(engine *xorm.EngineGroup, id int64) (*model.UserBasic, error) {
	if engine == nil {
		tool.Logger.Error(server.GetMsgByCode("en", server.InternalErrCode))
		return nil, errors.New(server.GetMsgByCode("en", server.InternalErrCode))
	}
	var entity model.UserBasic
	b, err := engine.Where("id=? AND deleted_at=0", id).Get(&entity)
	if err != nil {
		tool.Logger.Error(err)
		return nil, err
	}
	if !b {
		tool.Logger.Error("not found")
		return nil, nil
	}
	return &entity, nil
}
