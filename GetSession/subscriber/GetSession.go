package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	getSession "zufang/GetSession/proto/GetSession"
)

type GetSession struct{}

func (e *GetSession) Handle(ctx context.Context, msg *getSession.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *getSession.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
