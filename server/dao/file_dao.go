package dao

import (
	"cloud/server"
	"cloud/server/model"
	"cloud/server/vo"
	"cloud/tool"
	"errors"
	"github.com/go-xorm/xorm"
)

type FileDao struct {
}

// Add 插入单条记录
func (d FileDao) Add(engine *xorm.EngineGroup, entity *model.File) error {
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

// GetByID 根据id查询
func (d FileDao) GetByID(engine *xorm.EngineGroup, id string) (*model.File, error) {
	if engine == nil {
		tool.Logger.Error(server.GetMsgByCode("en", server.InternalErrCode))
		return nil, errors.New(server.GetMsgByCode("en", server.InternalErrCode))
	}
	var entity model.File
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

// EditByID 根据id编辑
func (d FileDao) EditByID(engine *xorm.EngineGroup, id string, item *model.File) error {
	if engine == nil || id == "" || item == nil {
		tool.Logger.Error(server.GetMsgByCode("en", server.InternalErrCode))
		return errors.New(server.GetMsgByCode("en", server.InternalErrCode))
	}

	time := tool.UnixSecond()
	item.UpdatedAt = &time

	_, err := engine.ID(id).AllCols().Update(item)
	if err != nil {
		tool.Logger.Error(err)
		return err
	}
	return nil
}

// GetList 分页查询列表
func (d FileDao) GetList(engine *xorm.EngineGroup, where string, whereParam []interface{}, page, pageSize int,
	sort map[string][]string) ([]vo.FileList, int64, error) {
	list := new([]vo.FileList)

	session := engine.Table("file").Where(where, whereParam...).Limit(pageSize, (page-1)*pageSize)
	if asc, ok := sort["asc"]; ok {
		session.Asc(asc...)
	}
	if desc, ok := sort["desc"]; ok {
		session.Desc(desc...)
	}
	count, err := session.FindAndCount(list)
	if err != nil {
		tool.Logger.Error(err)
		return nil, 0, err
	}

	return *list, count, nil
}
