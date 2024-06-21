package dao

import (
	"awesomeProject/entity"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RDB *redis.Client

func FindUser(UserPtr *entity.User) *gorm.DB {
	return DB.Where(UserPtr)
}

func CreateUser(user entity.User) *gorm.DB {
	return DB.Create(&user)
}

func DeleteUser(user entity.User) *gorm.DB {
	return DB.Delete(&user)
}

func UpdateUser(oldUserPtr *entity.User, newUserPtr *entity.User) *gorm.DB {
	return DB.Model(oldUserPtr).Updates(newUserPtr)
}
