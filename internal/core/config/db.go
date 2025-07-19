package config

import (
	"fmt"
	"time"

	"github.com/liuzhaomax/go-maxms/internal/core/ext"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type dbConfig struct {
	Type         string `mapstructure:"type"`
	Debug        bool   `mapstructure:"debug"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	Params       string `mapstructure:"params"`
	Endpoint     endpoint
}

func InitDB() (*gorm.DB, func(), error) {
	cfg.App.Logger.Info(ext.FormatInfo("数据库连接启动"))

	db, clean, err := cfg.Lib.DB.LoadDB()
	if err != nil {
		LogFailure(ext.ConnectionFailed, "数据库连接失败", err)

		return nil, clean, err
	}

	LogSuccess("数据库连接成功")

	return db, clean, err
}

func (d *dbConfig) LoadDB() (*gorm.DB, func(), error) {
	LogSuccess("数据库品种: " + d.Type)

	db, err := gorm.Open(mysql.Open(d.DSN()), &gorm.Config{
		Logger: InitGormLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	if d.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	clean := func() {
		err = sqlDB.Close()
		if err != nil {
			LogFailure(ext.ConnectionFailed, "数据库断开连接失败", err)
		}
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, clean, err
	}

	sqlDB.SetMaxIdleConns(d.MaxIdleConns)
	sqlDB.SetMaxOpenConns(d.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(d.MaxLifeTime) * time.Second)

	return db, clean, err
}

func (d *dbConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.Secret.Mysql.UserName, cfg.Secret.Mysql.PassWord, d.Endpoint.Host, d.Endpoint.Port, cfg.Secret.Mysql.Name, d.Params)
}
