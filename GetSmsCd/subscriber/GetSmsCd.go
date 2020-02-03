package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	getSmsCd "zufang/GetSmsCd/proto/GetSmsCd"
)

type GetSmsCd struct{}

func (e *GetSmsCd) Handle(ctx context.Context, msg *getSmsCd.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *getSmsCd.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
