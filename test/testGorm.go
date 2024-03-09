package main

import (
	"go_exam/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:sykj_2022@tcp(192.168.1.146:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	db.AutoMigrate(&models.UserBasic{})
	// Create
	user := &models.UserBasic{}
	user.Name = "申专"
	//user.HeartbeatTime = time.Now()
	//user.LogoutTime = time.Now()
	//user.LoginTime = time.Now()
	//db.Create(user)
	// Read
	//fmt.Println(db.First(user, 1)) // 根据整型主键查找
	//db.First(user, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	// Update - 将 product 的 price 更新为 200
	//db.Model(user).Update("PassWord", "1234")
	// Update - 更新多个字段
	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅 更新非零值字段
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	// Delete - 删除 product
	//db.Delete(&product, 1)
}
