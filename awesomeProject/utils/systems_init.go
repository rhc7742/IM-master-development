package utils

import (
	"awesomeProject/dao"
	"awesomeProject/entity"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitConfig() {
	viper.SetConfigName("app") //读取配置文件
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("初始化配置失败,err=" + err.Error())
	}
}

func InitMysql() {
	logger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	temp_db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{
		Logger: logger})
	if err != nil {
		panic("连接mysql 失败,err==" + err.Error())
	}
	dao.DB = temp_db
	dao.DB.AutoMigrate(&entity.User{})
	dao.DB.AutoMigrate(&entity.Message{})
	dao.DB.AutoMigrate(&entity.Contact{})
	dao.DB.AutoMigrate(&entity.GroupBasic{})
}

func InitRedis() {
	dao.RDB = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
}

const PublishKey = "websocket"

// 发布消息
func Publish(ctx context.Context, channel string, msg string) error {
	err := dao.RDB.Publish(ctx, channel, msg).Err()
	fmt.Println("publishing")

	return err
}

// 订阅消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := dao.RDB.Subscribe(ctx, channel)

	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("subscribing")
	if err != nil {
		panic(err)
	}
	return msg.Payload, err
}
