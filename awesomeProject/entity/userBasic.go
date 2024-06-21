package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UID        string `gorm:"uniqueIndex;size:15"`
	Name       string
	Password   string
	Phone      string `valid:"matches(/^1[3456789]\d{9}$/、/^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\d{8}$/)"`
	Email      string `valid:"email"`
	Identity   string
	RandomNum  string
	ClientIp   string
	ClientPort string
	IsLogOut   bool
	DeviceInfo string
}
