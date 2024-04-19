package core

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type DB struct {
	Type         string `mapstructure:"type"`
	Debug        bool   `mapstructure:"debug"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	Name         string `mapstructure:"name"`
	Params       string `mapstructure:"params"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Endpoint
}

func InitDB() (*gorm.DB, func(), error) {
	cfg.App.Logger.Info(FormatInfo("数据库连接启动"))
	db, clean, err := cfg.Lib.DB.LoadDB()
	if err != nil {
		LogFailure(ConnectionFailed, "数据库连接失败", err)
		return nil, clean, err
	}
	LogSuccess("数据库连接成功")
	return db, clean, err
}

func (d *DB) LoadDB() (*gorm.DB, func(), error) {
	LogSuccess(fmt.Sprintf("数据库品种: %s", d.Type))
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
			LogFailure(ConnectionFailed, "数据库断开连接失败", err)
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

func (d *DB) DSN() string {
	if d.Password == "" {
		d.Password = "123456"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		d.Username, d.Password, d.Endpoint.Host, d.Endpoint.Port, d.Name, d.Params)
}
