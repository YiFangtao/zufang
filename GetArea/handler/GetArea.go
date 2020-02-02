package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"time"
	"zufang/Ihomeweb/models"
	"zufang/Ihomeweb/pkg/logging"
	"zufang/Ihomeweb/utils"

	GETAREA "zufang/GetArea/proto/GetArea"
)

type GetArea struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetArea) GetArea(ctx context.Context, req *GETAREA.Request, rsp *GETAREA.Response) error {
	fmt.Println("请求地域信息 api/v1.0/areas")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//连接redis句柄
	redisConfigMap := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": utils.G_redis_pwd,
	}
	fmt.Println(redisConfigMap)
	//将map转成json
	redisConfig, _ := json.Marshal(redisConfigMap)
	//连接redis
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logging.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//获取缓存数据
	areasInfoValue := bm.Get("areas_info")
	//如果不为空,说明成功
	if areasInfoValue != nil {
		fmt.Println("获取到缓存发送给前端")
		//用来存放解码到json
		var areasInfo []map[string]interface{}
		//解码
		err = json.Unmarshal(areasInfoValue.([]byte), &areasInfo)
		//进行循环赋值
		for key, value := range areasInfo {
			fmt.Println(key, value)
			//创建对应数据类型进行赋值
			area := GETAREA.Response_Address{
				Aid:   int32(value["area_id"].(float64)),
				Aname: value["area_name"].(string),
			}
			rsp.Data = append(rsp.Data, &area)
		}
		return nil
	}

	fmt.Println("没有拿到缓存")

	//创建orm句柄
	o := orm.NewOrm()
	//接收地区信息到切片
	var areas []models.Area

	//查询全部地区
	qs := o.QueryTable("area")
	num, err := qs.All(&areas)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	if num == 0 {
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	logging.Info("写入缓存")

	//将数据写入缓存
	//将查询到到数据编码成json格式
	areasInfoStr, err := json.Marshal(areas)
	//放入缓存
	err = bm.Put("areas_info", areasInfoStr, time.Second*3600)
	if err != nil {
		fmt.Println("数据库数据存入缓存失败", err)
		rsp.Errno = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//返回地区信息
	for key, value := range areas {
		fmt.Println(key, value)
		area := GETAREA.Response_Address{
			Aid:   int32(value.Id),
			Aname: string(value.Name),
		}
		rsp.Data = append(rsp.Data, &area)
	}
	return nil
}
