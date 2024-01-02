package model

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var ModelUserSet = wire.NewSet(wire.Struct(new(ModelUser), "*"))

type ModelUser struct {
	DB *gorm.DB
}

func (m *ModelUser) QueryUserByUsername(c *gin.Context, username string, user *User) error {
	result := m.DB.WithContext(c).Where("username = ?", username).First(user)
	if result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}

func (m *ModelUser) QueryUserByUserID(c *gin.Context, userID string, user *User) error {
	result := m.DB.WithContext(c).Where("user_id = ?", userID).First(user)
	if result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}
