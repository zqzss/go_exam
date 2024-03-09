package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go_exam/docs"
	"go_exam/server"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any",ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/index", server.GetIndex)
	r.GET("/user/getUserList", server.GetUserList)
	r.GET("/user/createUser",server.CreateUser)
	r.GET("/user/deleteUser",server.DeleteUser)
	r.POST("/user/updateUser",server.UpdateUser)
	r.POST("/user/findUserByNameAndPwd",server.FindUserByNameAndPwd)
	r.GET("/user/sendMsg",server.SendMsg)
	return r
}
