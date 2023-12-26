package core

import (
	"context"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var TransactionSet = wire.NewSet(wire.Struct(new(Trans), "*"))

type Trans struct {
	DB *gorm.DB
}

func (t *Trans) NewTrans(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, Trans{}, tx)
}

func (t *Trans) GetTrans(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(Trans{})
	return v, v != nil
}

func (t *Trans) ExecTrans(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := t.GetTrans(ctx); ok {
		return fn(ctx)
	}
	return t.DB.Transaction(func(db *gorm.DB) error {
		return fn(t.NewTrans(ctx, db))
	})
}
