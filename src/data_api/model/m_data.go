package model

import (
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
		return result.Error
	}
	return nil
}
