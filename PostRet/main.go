package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"zufang/PostRet/handler"
	"zufang/PostRet/subscriber"

	PostRet "zufang/PostRet/proto/PostRet"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostRet"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	PostRet.RegisterPostRetHandler(service.Server(), new(handler.PostRet))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostRet", service.Server(), new(subscriber.PostRet))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostRet", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
