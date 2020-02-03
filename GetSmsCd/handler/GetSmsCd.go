package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"reflect"
	"time"
	"zufang/Ihomeweb/models"
	"zufang/Ihomeweb/pkg/logging"
	"zufang/Ihomeweb/utils"

	_ "github.com/astaxie/beego/cache/redis"

	"github.com/micro/go-micro/util/log"

	getSmsCd "zufang/GetSmsCd/proto/GetSmsCd"
)

type GetSmsCd struct{}

//获取短信验证码
func (e *GetSmsCd) GetSmsCd(ctx context.Context, req *getSmsCd.Request, rsp *getSmsCd.Response) error {
	logging.Info("Get smscd api")
	//初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*验证手机号是否正确*/
	//创建数据库orm句柄
	//使用手机号作为查询条件
	o := orm.NewOrm()
	user := models.User{Mobile: req.Mobile}
	err := o.Read(&user)
	//如果不报错 说明找得到
	//找得到就说明手机号已存在
	if err == nil {
		logging.Info("用户已存在")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*验证图片验证码是否正确*/
	//连接redis
	redisConfigMap := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": utils.G_redis_pwd,
	}
	logging.Info(redisConfigMap)

	//将map转化成json
	redisConfig, _ := json.Marshal(redisConfigMap)
	logging.Info(string(redisConfig))

	//连接redis数据库
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logging.Info("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logging.Info(req.Id, reflect.TypeOf(req.Id))

	//查询相关数据
	value := bm.Get(req.Id)
	if value == nil {
		logging.Info("redis获取失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logging.Info(value, reflect.TypeOf(value))

	//格式转换
	valueStr, err := redis.String(value, nil)
	logging.Info(valueStr, reflect.TypeOf(valueStr))

	//数据对比
	if valueStr != req.Text {
		logging.Info("图片验证码错误")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = "图片验证码错误"
		return nil
	}

	/*调用短信接口发送短信*/
	//创建随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := r.Intn(9999) + 1001
	logging.Info("验证码：", size)

	//发送短信的配置信息

	/*将短信验证码存入数据库中*/
	err = bm.Put(req.Mobile, size, time.Second*300)
	if err != nil {
		logging.Info("redis创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSmsCd) Stream(ctx context.Context, req *getSmsCd.StreamingRequest, stream getSmsCd.GetSmsCd_StreamStream) error {
	log.Logf("Received GetSmsCd.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&getSmsCd.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetSmsCd) PingPong(ctx context.Context, stream getSmsCd.GetSmsCd_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&getSmsCd.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
