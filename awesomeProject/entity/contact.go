package entity

import "gorm.io/gorm"

// 关系
type Contact struct {
	gorm.Model
	OwnerId  uint //谁的关系信息
	TargetId uint
	Type     int
	Desc     string
}
