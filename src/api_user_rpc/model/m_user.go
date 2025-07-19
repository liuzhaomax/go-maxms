package model

import (
	"context"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var ModelUserSet = wire.NewSet(wire.Struct(new(ModelUser), "*"))

type ModelUser struct {
	DB *gorm.DB
}

func (m *ModelUser) QueryUserByUserID(ctx context.Context, userID string, user *User) error {
	result := m.DB.WithContext(ctx).Where("user_id = ?", userID).First(user)
	if result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}
