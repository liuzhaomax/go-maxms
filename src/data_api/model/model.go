package model

import (
	"github.com/liuzhaomax/go-maxms/src/data_api/pb"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	Mobile string `gorm:"index:idx_mobile;unique;varchar(11);not null"`
}

func Model2PB(data *Data) *pb.DataRes {
	dataRes := &pb.DataRes{
		Id:     int32(data.ID),
		Mobile: data.Mobile,
	}
	return dataRes
}
