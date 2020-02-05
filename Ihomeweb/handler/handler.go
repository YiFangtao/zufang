package handler

import (
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"image"
	"image/png"
	"reflect"
	"regexp"
	"zufang/Ihomeweb/models"
	"zufang/Ihomeweb/pkg/logging"
	"zufang/Ihomeweb/utils"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/service/grpc"
	"net/http"
	DELETESESSION "zufang/DeleteSession/proto/DeleteSession"
	GETAREA "zufang/GetArea/proto/GetArea"
	GETIMAGECD "zufang/GetImageCd/proto/GetImageCd"
	GETSESSION "zufang/GetSession/proto/GetSession"
	GETSMSCD "zufang/GetSmsCd/proto/GetSmsCd"
	GETUSERINFO "zufang/GetUserInfo/proto/GetUserInfo"
	POSTAVATAR "zufang/PostAvatar/proto/PostAvatar"
	POSTLOGIN "zufang/PostLogin/proto/PostLogin"
	POSTRET "zufang/PostRet/proto/PostRet"
)

func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//创建新的grpc返回句柄
	server := grpc.NewService()
	//初始化服务
	server.Init()

	//创建获取地区的服务,并返回句柄
	getAreaService := GETAREA.NewGetAreaService("go.micro.srv.GetArea", server.Client())
	//调用函数并且获得返回数据
	rsp, err := getAreaService.GetArea(context.TODO(), &GETAREA.Request{})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//创建返回类型的切片
	var areaList []models.Area
	//循环读取服务返回的数据
	for _, value := range rsp.Data {
		area := models.Area{
			Id:     int(value.Aid),
			Name:   value.Aname,
			Houses: nil,
		}
		areaList = append(areaList, area)
	}

	//创建返回数据map
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   areaList,
	}

	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")

	//将返回的数据发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImageCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logging.Info("获取图形验证码 url")

	//创建服务
	server := grpc.NewService()
	server.Init()

	//连接服务
	getImageCdService := GETIMAGECD.NewGetImageCdService("go.micro.srv.GetImageCd", server.Client())
	//获取前端发送过来的唯一uuid
	logging.Info(ps.ByName("uuid"))
	rsp, err := getImageCdService.GetImageCd(context.TODO(), &GETIMAGECD.Request{
		Uuid: ps.ByName("uuid"),
	})

	//判断函数调用是否成功
	if err != nil {
		logging.Info(err)
		http.Error(w, err.Error(), 502)
		return
	}

	//处理前端发送过来的图片信息
	var img image.RGBA
	img.Stride = int(rsp.Stride)
	img.Rect.Min.X = int(rsp.Min.X)
	img.Rect.Min.Y = int(rsp.Min.Y)
	img.Rect.Max.X = int(rsp.Max.X)
	img.Rect.Max.Y = int(rsp.Max.Y)
	img.Pix = rsp.Pix

	var image captcha.Image
	image.RGBA = &img

	//将图片发送给前端
	png.Encode(w, image)
}

func GetSmsCd(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logging.Info("获取短信验证码 ")

	//创建服务
	service := grpc.NewService()
	service.Init()

	//获取前端发送过来的手机号
	mobile := ps.ByName("mobile")
	logging.Info(mobile)

	//后端进行正则匹配
	myreg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	bo := myreg.MatchString(mobile)

	//如果手机号错误
	if bo == false {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "手机号错误",
		}
		//设置返回数据格式
		w.Header().Set("Content-Type", "application/json")
		//将错误发送给前端
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logging.Info(err)
			return
		}
		logging.Info("手机号错误返回")
		return
	}

	//获取url携带的验证码
	logging.Info(r.URL.Query())
	//获取url携带的参数
	text := r.URL.Query()["text"][0]
	id := r.URL.Query()["id"][0]

	//调用服务
	smsCdService := GETSMSCD.NewGetSmsCdService("go.micro.srv.GetSmsCd", service.Client())
	rsp, err := smsCdService.GetSmsCd(context.TODO(), &GETSMSCD.Request{
		Mobile: mobile,
		Id:     id,
		Text:   text,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logging.Info(err)
		return
	}
	//创建返回map
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//设置返回格式
	w.Header().Set("Content-Type", "application/json")

	//将数据返回给前端
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logging.Info(err)
		return
	}

}

func PostReg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Info("注册请求 ")
	/*获取前端发送过来的json数据*/
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for key, value := range request {
		logging.Info(key, value, reflect.TypeOf(value))
	}

	//数据校验
	if request["mobile"].(string) == "" || request["password"].(string) == "" || request["sms_code"].(string) == "" {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_NODATA,
			"errmsg": "信息有误 请重新输入",
		}
		//设置番薯数据的格式
		w.Header().Set("Content-Type", "application/json")
		//发送给前端
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logging.Info(err)
			return
		}
		return
	}

	//创建服务
	service := grpc.NewService()
	service.Init()

	//连接服务 将数据发送给注册服务
	postRetService := POSTRET.NewPostRetService("go.micro.srv.PostRet", service.Client())
	rsp, err := postRetService.PostRet(context.TODO(), &POSTRET.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
		SmsCode:  request["sms_code"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logging.Debug(err)
		return
	}

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}

	//读取cookie
	cookie, err := r.Cookie("userLogin")
	//如果读取失败 或者cookie的value不存在 则创建cookie
	if err != nil || "" == cookie.Value {
		//创建1个cookie对象
		cookie := http.Cookie{
			Name:   "userLogin",
			Value:  rsp.SessionID,
			Path:   "/",
			MaxAge: 600,
		}
		//对浏览器的cookie进行设置
		http.SetCookie(w, &cookie)
	}

	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	//将数据发送给前端
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logging.Debug(err)
		return
	}
	return
}

func GetSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Info("获取Session url")

	//获取cookie
	cookie, err := r.Cookie("userLogin")
	if err != nil {
		//直接返回 说明用户未登陆
		response := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		//设置返回数据的格式
		w.Header().Set("Content-Type", "application/json")
		//将数据发给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	//创建服务
	service := grpc.NewService()
	service.Init()

	//创建句柄
	getSessionService := GETSESSION.NewGetSessionService("go.micro.srv.GetSession", service.Client())
	rsp, err := getSessionService.GetSession(context.TODO(), &GETSESSION.Request{
		SessionID: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//将获取到的用户名返回给前端
	data := make(map[string]string)
	data["name"] = rsp.Username
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	//设置返回数据的格式
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Info("登陆 api/v1.0/sessions")
	//获取前端post请求的内容
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for key, value := range request {
		logging.Info(key, value, reflect.TypeOf(value))
	}

	//判断账号密码是否为空
	if request["mobile"] == "" || request["password"] == "" {
		response := map[string]interface{}{
			"errno":  utils.RECODE_DBERR,
			"errmsg": "用户名或密码不能为空",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 503)
			logging.Debug(err)
			return
		}
		return
	}

	//创建连接
	service := grpc.NewService()
	service.Init()

	loginService := POSTLOGIN.NewPostLoginService("go.micro.srv.PostLogin", service.Client())
	rsp, err := loginService.PostLogin(context.TODO(), &POSTLOGIN.Request{
		Mobile:   request["mobile"].(string),
		Password: request["password"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logging.Debug(err)
		return
	}

	/*设置cookie*/
	//读取cookie
	cookie, err := r.Cookie("userLogin")
	if err != nil || cookie.Value == "" {
		cookie := http.Cookie{
			Name:   "userLogin",
			Value:  rsp.SessionID,
			Path:   "/",
			MaxAge: 600,
		}
		http.SetCookie(w, &cookie)
	}

	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回的数据格式
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func DeleteSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Info("DELETE 退出登陆 /api/v1.0/session")

	//获取cookie
	cookie, err := r.Cookie("userLogin")
	if err != nil || cookie.Value == "" {
		//准备回传数据
		response := map[string]interface{}{
			"errno":  utils.RECODE_DATAERR,
			"errmsg": utils.RecodeText(utils.RECODE_DATAERR),
		}
		//设置返回的数据格式
		w.Header().Set("Content-Type", "application/json")
		//发送给前端
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	//创建服务
	service := grpc.NewService()
	service.Init()
	//调用服务
	deleteSessionService := DELETESESSION.NewDeleteSessionService("go.micro.srv.DeleteSession", service.Client())

	rsp, err := deleteSessionService.DeleteSession(context.TODO(), &DELETESESSION.Request{
		SessionID: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//删除sessionID
	cookie, err = r.Cookie("userLogin")
	if cookie.Value != "" || err == nil {
		cookie := http.Cookie{
			Name:   "userLogin",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, &cookie)
	}

	//返回数据
	response := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
	}
	//设置返回数据的格式
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Info("GetUserInfo 获取用户信息")
	//初始化服务
	service := grpc.NewService()
	service.Init()

	//创建句柄
	userInfoService := GETUSERINFO.NewGetUserInfoService("go.micro.srv.GetUserInfo", service.Client())
	//获取用户的登陆信息
	cookie, err := r.Cookie("userLogin")
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logging.Debug(err)
			return
		}
	}

	//将信息发送给前端
	rsp, err := userInfoService.GetUserInfo(context.TODO(), &GETUSERINFO.Request{
		SessionID: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logging.Debug(err)
		return
	}

	//准备1哥数据的map
	data := make(map[string]interface{})
	//将数据发送给前端
	data["user_id"] = rsp.UserID
	data["name"] = rsp.Username
	data["mobile"] = rsp.Mobile
	data["real_name"] = rsp.RealName
	data["id_card"] = rsp.IDCard
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	//设置格式
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logging.Debug(err)
		return
	}
	return
}

func PostAvatar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logging.Info("上传用户头像")

	//创建服务
	service := grpc.NewService()
	service.Init()

	//创建句柄
	postAvatarService := POSTAVATAR.NewPostAvatarService("go.micro.srv.PostAvatar", service.Client())

	//查看登陆信息
	cookie, err := r.Cookie("userLogin")
	//如果没有登陆就返回错误
	if err != nil {
		resp := map[string]interface{}{
			"errno":  utils.RECODE_SESSIONERR,
			"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logging.Debug(err)
			return
		}
		return
	}

	//接收前端发送过来的图片
	file, header, err := r.FormFile("avatar")
	//判断是否接收成功
	if err != nil {
		logging.Debug("PostAvatar c.GetFile(avatar) err", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logging.Debug(err)
			return
		}
		return
	}
	//打印基本信息
	logging.Info(file, header)
	logging.Info("文件大小", header.Size)
	logging.Info("文件名", header.Filename)

	//二进制的空间用来存储文件
	fileBuffer := make([]byte, header.Size)
	//将文件读取到fileBuffer里
	_, err = file.Read(fileBuffer)
	if err != nil {
		logging.Debug("PostAvatar c.GetFile(avatar) err", err)
		resp := map[string]interface{}{
			"errno":  utils.RECODE_IOERR,
			"errmsg": utils.RecodeText(utils.RECODE_IOERR),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), 503)
			logging.Debug(err)
			return
		}
		return
	}

	//调用函数传入数据
	rsp, err := postAvatarService.PostAvatar(context.TODO(), &POSTAVATAR.Request{
		Avatar:    fileBuffer,
		SessionID: cookie.Value,
		FileSize:  header.Size,
		FileName:  header.Filename,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		logging.Debug(err)
		return
	}

	//准备回传数据
	data := make(map[string]interface{})
	//url拼接后回传数据
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)
	resp := map[string]interface{}{
		"errno":  rsp.Errno,
		"errmsg": rsp.Errmsg,
		"data":   data,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), 503)
		logging.Debug(err)
		return
	}
	return
}
