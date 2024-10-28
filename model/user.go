package model

import (
	"gorm.io/gorm"
	"learn/internal"
)

func init() {
	internal.DB.AutoMigrate(&User{})
}

type User struct {
	gorm.Model
	Name string
	Coin int
}
