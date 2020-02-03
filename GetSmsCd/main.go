package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"zufang/GetSmsCd/handler"
	"zufang/GetSmsCd/subscriber"

	GetSmsCd "zufang/GetSmsCd/proto/GetSmsCd"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetSmsCd"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetSmsCd.RegisterGetSmsCdHandler(service.Server(), new(handler.GetSmsCd))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSmsCd", service.Server(), new(subscriber.GetSmsCd))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSmsCd", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
