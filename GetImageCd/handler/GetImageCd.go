package handler

import (
	"context"
	"encoding/json"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/micro/go-micro/util/log"
	"image/color"
	"time"
	"zufang/Ihomeweb/pkg/logging"
	"zufang/Ihomeweb/utils"

	getImageCd "zufang/GetImageCd/proto/GetImageCd"
)

type GetImageCd struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetImageCd) GetImageCd(ctx context.Context, req *getImageCd.Request, rsp *getImageCd.Response) error {
	logging.Info("-------GET /api/v1.0/imagecd")
	//创建1个句柄
	cap := captcha.New()
	//通过句柄调用 字体文件
	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}

	//设置图片大小
	cap.SetSize(128, 64)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})

	//生成图片
	img, str := cap.Create(4, captcha.NUM)
	logging.Info(str)

	//解引用
	b := *img
	c := *(b.RGBA)
	//成功返回
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	//图片信息
	rsp.Pix = c.Pix
	rsp.Stride = int64(c.Stride)
	rsp.Max = &getImageCd.Response_Point{
		X: int64(c.Rect.Max.X),
		Y: int64(c.Rect.Max.Y),
	}
	rsp.Min = &getImageCd.Response_Point{
		X: int64(c.Rect.Min.X),
		Y: int64(c.Rect.Min.Y),
	}

	//将uuid与随机数验证码对应存储在redis缓存中
	redisConfigMap := map[string]string{
		"key":      "ihome",
		"conn":     utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum":    utils.G_redis_dbnum,
		"password": utils.G_redis_pwd,
	}
	logging.Info(redisConfigMap)
	redisConfig, _ := json.Marshal(redisConfigMap)

	//连接数据库 创建句柄
	bm, err := cache.NewCache("redis", string(redisConfig))
	if err != nil {
		logging.Info("GetImage() cache.NewCache err", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}

	//验证码进行1小时缓存
	bm.Put(req.Uuid, str, 300*time.Second)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetImageCd) Stream(ctx context.Context, req *getImageCd.StreamingRequest, stream getImageCd.GetImageCd_StreamStream) error {
	log.Logf("Received GetImageCd.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&getImageCd.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetImageCd) PingPong(ctx context.Context, stream getImageCd.GetImageCd_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&getImageCd.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
