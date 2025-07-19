package db_auto_migrate_user_test

import (
	"testing"

	"github.com/liuzhaomax/go-maxms/internal/core/config"
	"github.com/liuzhaomax/go-maxms/internal/core/ext"
)

func TestAutoMigrate(t *testing.T) {
	err := AutoMigrate()
	if err != nil {
		config.LogFailure(ext.Unknown, "数据库表创建失败", err)

		return
	}
}
