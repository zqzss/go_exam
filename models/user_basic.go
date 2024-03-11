package models

import (
	"fmt"
	"go_exam/utils"
	"gorm.io/gorm"
	"time"
)


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
