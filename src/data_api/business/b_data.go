package business

import (
	"context"
	"github.com/google/wire"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/data_api/model"
	"github.com/liuzhaomax/go-maxms/src/data_api/pb"
	"github.com/sirupsen/logrus"
)

var BusinessDataSet = wire.NewSet(wire.Struct(new(BusinessData), "*"))

type BusinessData struct {
	Model  *model.ModelData
	Logger *logrus.Logger
	Tx     *core.Trans
}

func (b *BusinessData) GetDataById(ctx context.Context, req *pb.IdRequest) (*pb.DataRes, error) {
	var data *model.Data
	err := b.Tx.ExecTrans(ctx, func(ctx context.Context) error {
		err := b.Model.QueryDataById(req, data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		b.Logger.WithField("失败方法", core.GetFuncName()).Info(core.FormatError(core.NotFound, "事务执行失败", err))
		return nil, err
	}
	res := model.Model2PB(data)
	b.Logger.WithField("成功方法", core.GetFuncName()).Info(core.FormatInfo("响应成功"))
	return res, nil
}
