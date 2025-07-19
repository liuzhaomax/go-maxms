package db_auto_migrate_user

import (
	"fmt"
	"github.com/liuzhaomax/go-maxms/internal/core"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

func AutoMigrate() error {
	cfg := core.GetConfig()
	db, err := gorm.Open(mysql.Open(DSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		core.LogFailure(core.ConnectionFailed, "数据库连接失败", err)
		return nil
	}
	dbType := strings.ToLower(cfg.Lib.DB.Type)
	if dbType == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	err = db.AutoMigrate(new(model.User))
	if err != nil {
		return err
	}
	createAdmin(db) // 添加一条数据
	return nil
}

func createAdmin(db *gorm.DB) {
	cfg := core.GetConfig()
	user := &model.User{}
	result := db.First(user)
	salt, encodedPwd := core.GetEncodedPwd("admin")
	cfg.App.Salt = salt
	// 将salt更新到vault
	if cfg.App.Enabled.Vault {
		cfg.PutSalt()
	}
	if result.RowsAffected == 0 {
		user.UserID = core.ShortUUID()
		user.Username = "admin"
		user.Password = encodedPwd
		user.Mobile = "+8613012345678"
		user.Email = "admin@maxblog.cn"
		res := db.Create(&user)
		if res.RowsAffected == 0 {
			core.LogFailure(core.DBDenied, "admin创建失败", res.Error)
			panic(res.Error)
		}
	} else {
		res := db.Model(user).Where("user_id = ?", user.UserID).Update("password", encodedPwd)
		if res.RowsAffected == 0 {
			core.LogFailure(core.DBDenied, "admin更新失败", res.Error)
			panic(res.Error)
		}
	}
}

func DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		"root", "123456", "127.0.0.1", "3306", "go_maxms", "charset=utf8mb4&parseTime=True&loc=Local")
}
