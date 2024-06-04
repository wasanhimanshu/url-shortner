package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlStorage() (*gorm.DB, error) {
	dsnString := Envs.DBUser + ":" + Envs.DBPasswd + "@tcp(" + Envs.DBHost + ":" + Envs.DBPort + ")/" + Envs.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsnString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{}, &Url{})
	return db, nil
}
