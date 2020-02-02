package handler

import (
	"context"
	"encoding/json"
	"zufang/Ihomeweb/models"

	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro/service/grpc"
	"net/http"
	GETAREA "zufang/GetArea/proto/GetArea"
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

	w.Header().Set("Content-Type", "application/json")

	//将返回的数据发送给前端
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
