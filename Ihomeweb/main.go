package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"net/http"
	"zufang/Ihomeweb/handler"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.Ihomeweb"),
		web.Version("latest"),
		web.Address(":8999"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	//获取地区信息
	router.GET("/api/v1.0/areas", handler.GetArea)
	//获取图形验证码
	router.GET("/api/v1.0/imagecode/:uuid", handler.GetImageCd)
	//获取短信验证码
	router.GET("/api/v1.0/smscode/:mobile", handler.GetSmsCd)
	//注册
	router.POST("/api/v1.0/users", handler.PostReg)
	//获取session
	router.GET("/api/v1.0/session", handler.GetSession)
	//登陆
	router.POST("/api/v1.0/sessions", handler.PostLogin)
	//退出登陆
	router.DELETE("/api/v1.0/session", handler.DeleteSession)
	//获得用户信息
	router.GET("/api/v1.0/user", handler.GetUserInfo)
	//上传头像
	router.POST("/api/v1.0/user/avatar", handler.PostAvatar)

	//映射静态页面
	router.NotFound = http.FileServer(http.Dir("html"))
	//注册服务
	service.Handle("/", router)

	//运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
