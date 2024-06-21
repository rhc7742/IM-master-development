package entity

import "gorm.io/gorm"

// 群
type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Desc    string
	Type    string //冲钱加人
}
