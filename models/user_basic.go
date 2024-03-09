package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go_exam/utils"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"time"
)

type Node struct{
	Conn *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

type UserBasic struct {
	gorm.Model
	Name       string
	Password   string
	Phone      string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email      string `valid:"email"`
	Identity   string
	ClientIP   string
	ClientPort string
	Salt       string
	//LoginTime     time.Time
	//HeartbeatTime time.Time
	//LogoutTime time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout   bool
	DeviceInfo string
}
func (table *UserBasic) TableName() string {
	return "user_basic"
}

//消息
type Message struct {
	gorm.Model
	FormId uint
	//发送者
	TargetId uint
	//接受者
	Type string //消息类型 群聊 私聊 广播
	Media int //消息类型 文字 图片 音频
	Content string //消息内容
	Pic string
	Url string
	Desc string
	Amount int //其他数字统计
}
func (table *Message) TableName() string {
	return "message"
}

//群信息
type GroupBasic struct {

	//发送消息 接受消息 需要 ：发送者ID ，接受者ID ，消息类型，发送的内容，发送类型
	gorm.Model
	Name string
	OwnerId uint
	Icon string
	Type int
	Desc string
}

func (table *GroupBasic) TableName() string{
	return "group_basic"
}

//人员关系
type Contact struct {
	gorm.Model
	OwnerId uint //谁的关系信息
	TargetId uint //对应的谁
	Type int //对应的类型  0  1  3
	Desc string
}
func (table *Contact) TableName() string {
	return "contact"
}

func GetUserList() []UserBasic {
	data := make([]UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user UserBasic) {
	utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) {
	utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) {

	utils.DB.Model(&user).Updates(user)
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

func FindUserByNameAndPwd(name , password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, password).First(&user)

	// token加密
	str := fmt.Sprintf("%d",time.Now().Unix())
	temp := utils.Md5Encode(str)
	utils.DB.Model(&user).Where("id = ?",user.ID).Update("identity",temp)

	return user
}
