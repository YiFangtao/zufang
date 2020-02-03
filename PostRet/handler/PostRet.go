package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"time"
	"zufang/Ihomeweb/models"
	"zufang/Ihomeweb/pkg/logging"
	"zufang/Ihomeweb/utils"

	"github.com/micro/go-micro/util/log"

	postRet "zufang/PostRet/proto/PostRet"
)

type PostRet struct{}

func GetMd5String(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostRet) PostRet(ctx context.Context, req *postRet.Request, rsp *postRet.Response) error {
	logging.Info("PostRet 注册")
	//初始化错误码
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	/*验证短信验证码*/

	/*redis操作*/
	//配置缓存参数
	redisConfigMap := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": utils.G_redis_pwd,
	}
	//将map转成json
	redisConfig, _ := json.Marshal(redisConfigMap)

	//创建redis句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logging.Debug("缓存创建失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//通过手机号获取到短信验证码
	smsCode := bm.Get(req.Mobile)
	if smsCode == nil {
		logging.Debug("获取数据失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//短信验证码对比
	smsCodeStr, _ := redis.String(smsCode, nil)
	if smsCodeStr != req.SmsCode {
		logging.Debug("短信验证码错误")
		rsp.Errno = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*将数据存入数据库*/
	o := orm.NewOrm()
	user := models.User{
		Name:         req.Mobile,
		PasswordHash: GetMd5String(req.Password),
		Mobile:       req.Mobile,
	}
	id, err := o.Insert(&user)
	if err != nil {
		logging.Debug("注册数据失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	logging.Info("user_id", id)

	/*创建sessionID 保证唯一性*/
	sessionID := GetMd5String(req.Mobile + req.Password)
	//返回给客户端
	rsp.SessionID = sessionID

	/*以sessionID为key的一部分创建session*/
	bm.Put(sessionID+"name", user.Mobile, time.Second*3600)
	bm.Put(sessionID+"user_id", id, time.Second*3600)
	bm.Put(sessionID+"mobile", user.Mobile, time.Second*3600)

	return nil

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostRet) Stream(ctx context.Context, req *postRet.StreamingRequest, stream postRet.PostRet_StreamStream) error {
	log.Logf("Received PostRet.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&postRet.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostRet) PingPong(ctx context.Context, stream postRet.PostRet_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&postRet.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
