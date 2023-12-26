package core

import (
	"fmt"
	"github.com/liuzhaomax/go-maxms-template/src/data_api/model"
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
		cfg.App.Logger.WithField("失败方法", GetFuncName()).Fatal(FormatError(Unknown, "数据库连接失败", err))
		return nil, clean, err
	}
	err = cfg.AutoMigrate(db)
	if err != nil {
		cfg.App.Logger.WithField("失败方法", GetFuncName()).Fatal(FormatError(Unknown, "数据库自动迁移失败", err))
		return nil, clean, err
	}
	cfg.App.Logger.Info(FormatInfo("数据库连接成功"))
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
	cfg.App.Logger.Info(FormatInfo(fmt.Sprintf("数据库品种: %s", cfg.DB.Type)))
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
			cfg.App.Logger.WithField("失败方法", GetFuncName()).Error(FormatError(Unknown, "数据库断开连接失败", err))
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
	err := db.AutoMigrate(new(model.Data))
	if err != nil {
		return err
	}
	createAdmin(db)
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
	var data model.Data
	db.First(&data)
	if data.ID != 1 {
		data.Mobile = "130123456789"
		db.Create(&data)
	}
}
