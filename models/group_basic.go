package models

import "gorm.io/gorm"

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
