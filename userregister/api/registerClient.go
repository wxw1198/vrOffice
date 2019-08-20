package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"

	"github.com/wxw1198/vrOffice/userregister/api/handler"
	"github.com/wxw1198/vrOffice/userregister/api/config"
	"github.com/wxw1198/vrOffice/userregister/proto"
)

func main() {
	//获取相关配置
	if !config.ParseConfig(&config.DefaultConfig) {
		log.Fatal("parse config err")
		return
	}

	//用于服务发现的地址
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.DefaultConfig.Etcdv3,
		}
	})

	service := micro.NewService(
		micro.Registry(reg),
		micro.Name(config.DefaultConfig.MicroName),
		micro.Address(config.DefaultConfig.ListenLocalAddr),
	)

	// Init to parse flags
	service.Init()

	//设置服务器的公网地址
	service.Server().Init(server.Advertise(config.DefaultConfig.AdvertiseAddr))

	//指定对urlgo.micro.api/register处理的handler
	//proto.RegisterRegisterHandler(service.Server(), &handler.RegisterHandler{})
	proto.RegisterRegisterHandler(service.Server(), &handler.RegisterHandler{
		//Create Service Client
		Client: proto.NewRegisterService("go.micro.srv.register", service.Client()),
	})

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
