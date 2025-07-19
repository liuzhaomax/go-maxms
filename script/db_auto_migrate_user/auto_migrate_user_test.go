package db_auto_migrate_user

import (
	"github.com/liuzhaomax/go-maxms/internal/core"
	"testing"
)

func TestAutoMigrate(t *testing.T) {
	err := AutoMigrate()
	if err != nil {
		core.LogFailure(core.Unknown, "数据库表创建失败", err)
		return
	}
}
