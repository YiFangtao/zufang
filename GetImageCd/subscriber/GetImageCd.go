package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	getImageCd "zufang/GetImageCd/proto/GetImageCd"
)

type GetImageCd struct{}

func (e *GetImageCd) Handle(ctx context.Context, msg *getImageCd.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *getImageCd.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
