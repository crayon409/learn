package internal

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	if DB, err = gorm.Open(mysql.Open("root:admin@(localhost:30002)/learn?parseTime=true&loc=Local"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}); err != nil {
		panic(err)
	}
}
