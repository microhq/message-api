package main

import (
	"log"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/message-api/handler"
	proto "github.com/micro/message-api/proto/message"
	proto2 "github.com/micro/message-srv/proto/message"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.message"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Init()

	proto.RegisterMessageHandler(service.Server(), new(handler.Message))

	handler.Client = proto2.NewMessageClient("go.micro.srv.message", service.Client())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
