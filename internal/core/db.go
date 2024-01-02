package core

import (
	"fmt"
	"github.com/liuzhaomax/go-maxms/src/api_user/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strings"
	"time"
)

func InitDB() (*gorm.DB, func(), error) {
	cfg.App.Logger.Info(FormatInfo("数据库连接启动"))
	db, clean, err := cfg.LoadDB()
	if err != nil {
		LogFailure(ConnectionFailed, "数据库连接失败", err)
		return nil, clean, err
	}
	err = cfg.AutoMigrate(db)
	if err != nil {
		LogFailure(Unknown, "数据库表创建失败", err)
		return nil, clean, err
	}
	LogSuccess("数据库连接成功")
	createAdmin(db) // 添加一条数据
	return db, clean, err
}

func (cfg *Config) LoadDB() (*gorm.DB, func(), error) {
	gormLogger := logger.New(
		cfg.App.Logger,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			ParameterizedQueries:      true,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	LogSuccess(fmt.Sprintf("数据库品种: %s", cfg.DB.Type))
	db, err := gorm.Open(mysql.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	if cfg.DB.Debug {
		db = db.Debug()
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	clean := func() {
		err = sqlDB.Close()
		if err != nil {
			LogFailure(ConnectionFailed, "数据库断开连接失败", err)
		}
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, clean, err
	}
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DB.MaxLifeTime) * time.Second)
	return db, clean, err
}

func (cfg *Config) AutoMigrate(db *gorm.DB) error {
	dbType := strings.ToLower(cfg.DB.Type)
	if dbType == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	err := db.AutoMigrate(new(model.User))
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DSN() string {
	if db.Password == "" {
		db.Password = "123456"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		db.Username, db.Password, db.Host, db.Port, db.Name, db.Params)
}

func createAdmin(db *gorm.DB) {
	user := &model.User{}
	result := db.First(user)
	salt, encodedPwd := GetEncodedPwd("admin")
	cfg.App.Salt = salt
	// 将salt更新到vault
	if cfg.App.Enabled.Vault {
		cfg.PutSalt()
	}
	if result.RowsAffected == 0 {
		user.UserID = ShortUUID()
		user.Username = "admin"
		user.Password = encodedPwd
		user.Mobile = "+8613012345678"
		user.Email = "admin@maxblog.cn"
		res := db.Create(&user)
		if res.RowsAffected == 0 {
			LogFailure(DBDenied, "admin创建失败", res.Error)
			panic(res.Error)
		}
	} else {
		res := db.Model(user).Where("user_id = ?", user.UserID).Update("password", encodedPwd)
		if res.RowsAffected == 0 {
			LogFailure(DBDenied, "admin更新失败", res.Error)
			panic(res.Error)
		}
	}
}
