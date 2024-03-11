package server

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_exam/models"
	"go_exam/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}

func MsgHandler(c *gin.Context, ws *websocket.Conn) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		msgType,p,_ := ws.ReadMessage()
		fmt.Printf("收到消息: ",string(p))
		err = ws.WriteMessage(msgType, append([]byte(m),p...))
		if err != nil {
			log.Fatal(err)
		}
	}
}

// GetUserList
// @Tags 用户模块
// @Summary 所有用户
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{"message": data})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	email := c.Query("email")
	phone := c.Query("phone")
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(200, gin.H{"code": -1, "data": "", "message": "用户名已注册！"})
		return
	}

	data = models.FindUserByEmail(email)
	if data.Email != "" {
		c.JSON(200, gin.H{"code": -1, "data": "", "message": "邮箱已注册！"})
		return
	}

	data = models.FindUserByPhone(phone)
	if data.Phone != "" {
		c.JSON(200, gin.H{"code": -1, "data": "", "message": "手机号已注册！"})
		return
	}

	if password != repassword {
		c.JSON(200, gin.H{"code": -1, "data": "", "message": "两次密码不一致！"})
		return
	}
	user.Email = email
	user.Phone = phone
	salt := "abc123" + user.Name
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt)
	models.CreateUser(user)
	c.JSON(200, gin.H{"code": 0, "data": "", "message": "新增用户成功！"})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{"code": 0, "data": "", "message": "删除用户成功"})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	fmt.Printf("update: %+v \n", user)
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"code": -1, "data": "", "message": "修改参数不匹配！"})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{"code": 0, "data": "", "message": "修改用户成功！"})
	}
}

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{"message": "该用户不存在！"})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(200, gin.H{"message": "密码不正确"})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{"code": 0, "data": data, "message": "修改用户成功！"})
}
