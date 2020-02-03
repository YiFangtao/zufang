package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/gomodule/redigo/redis"
	"zufang/Ihomeweb/pkg/logging"
	"zufang/Ihomeweb/utils"

	"github.com/micro/go-micro/util/log"

	getSession "zufang/GetSession/proto/GetSession"
)

type GetSession struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetSession) GetSession(ctx context.Context, req *getSession.Request, rsp *getSession.Response) error {
	logging.Info("获取session信息")
	//初始化返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//构建连接缓存的数据
	redisConfigMap := map[string]string{
		"key":      utils.G_server_name,
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": utils.G_redis_pwd,
	}
	//将map转化成json
	redisConfig, _ := json.Marshal(redisConfigMap)
	//创建redis句柄 连接数据库
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logging.Debug("redis连接失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	//从缓存中获取session
	username := bm.Get(req.SessionID + "name")
	//如果不存在 返回失败
	if username == nil {
		logging.Debug("数据不存在", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}
	//如果存在 就返回
	rsp.Username, _ = redis.String(username, nil)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSession) Stream(ctx context.Context, req *getSession.StreamingRequest, stream getSession.GetSession_StreamStream) error {
	log.Logf("Received GetSession.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&getSession.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetSession) PingPong(ctx context.Context, stream getSession.GetSession_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&getSession.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
