package userModel

import (
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var (
	db *gorm.DB
	err error
)

type User struct {
	ID uint `gorm:"primary_key"`
	Name string `gorm:"not null"`
	Age int `gorm:"not null;size:11"`
	Strength int `gorm:"default:10;size:10"`
	Email string `gorm:"not null;unique"`
	Price float64 `gorm:"not null;type:decimal(10,2)"`
	Image string
	Description string `gorm:"type:longtext"`
	DeletedAt *time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

func init() {
	db, err = gorm.Open("mysql","root:@(localhost:3306)/loggerman?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	if !db.HasTable(&User{}) {
		db.AutoMigrate(&User{})
	}
}

func Connection() *gorm.DB {
	return db
}

