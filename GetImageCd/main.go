package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"zufang/GetImageCd/handler"
	"zufang/GetImageCd/subscriber"

	GetImageCd "zufang/GetImageCd/proto/GetImageCd"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetImageCd"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetImageCd.RegisterGetImageCdHandler(service.Server(), new(handler.GetImageCd))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetImageCd", service.Server(), new(subscriber.GetImageCd))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetImageCd", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
