package model

import (
	"errors"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/src/data_api/pb"
	"gorm.io/gorm"
)

var ModelDataSet = wire.NewSet(wire.Struct(new(ModelData), "*"))

type ModelData struct {
	DB *gorm.DB
}

func (m *ModelData) QueryDataById(req *pb.IdRequest, data *Data) error {
	result := m.DB.First(data, req.Id)
	if result.RowsAffected == 0 {
		return errors.New("所查询的数据不存在")
	}
	return nil
}
