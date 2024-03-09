package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

var (
	DB *gorm.DB
	Red *redis.Client
)
const (
	PublishKey = "websocket"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config  app inited 。。。。")
}
func InitMySQL() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		})
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")),
		&gorm.Config{Logger: newLogger})
	fmt.Println(" MySQL inited 。。。。")
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{Addr: viper.GetString("redis.addr"), Password: viper.GetString("password"), DB: viper.GetInt("redis.DB"), PoolSize: viper.GetInt("redis.poolSize"), MinIdleConns: viper.GetInt("redis.minIdleConn")})
	pong,err := Red.Ping(Red.Context()).Result()
	if err != nil{
		fmt.Println("init redis 。。。",err)
	}else {
		fmt.Println("Redis inited 。。。",pong)
	}
}

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	temStr := h.Sum(nil)
	return hex.EncodeToString(temStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt string, password string) bool {
	md := Md5Encode(plainpwd + salt)
	fmt.Println(md + "           " + password)
	return md == password
}

// Publish 发布消息到Redis
func Publish(ctx context.Context,channel string,msg string) error{
	var err error
	fmt.Println("Publish 。。。",msg)
	err = Red.Publish(ctx,channel,msg).Err()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

//Subscribe 订阅redis消息
func Subscribe(ctx context.Context,channel string)(string,error){
	sub := Red.Subscribe(ctx,channel)
	fmt.Println("Subscribe 。。。",ctx)
	msg,err := sub.ReceiveMessage(ctx)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("Subscribe 。。。。", msg.Payload)
	return msg.Payload,err
}