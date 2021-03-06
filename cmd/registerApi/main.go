package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/wxw1198/vrOffice/log"

	"github.com/wxw1198/vrOffice/userbaseoperation/api/config"
	"github.com/wxw1198/vrOffice/userbaseoperation/api/handler"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
)

func main() {
	//获取相关配置
	if config.ParseConfig(&config.DefaultConfig) == false {
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
		//micro.name response proto/userBaseOperation.proto service name
		micro.Name(config.DefaultConfig.MicroName),
		micro.Address(config.DefaultConfig.ListenLocalAddr),
	)

	// Init to parse flags
	service.Init()

	//设置服务器的公网地址
	//service.Server().Init(server.Advertise(config.DefaultConfig.AdvertiseAddr))

	//指定对urlgo.micro.api/register处理的handler
	//proto.RegisterRegisterHandler(service.Server(), &handler.UserBaseOperationHandler{})
	proto.RegisterUserBaseOpsHandler(service.Server(), &handler.UserBaseOperationHandler{
		//Create Service Client
		Client: proto.NewUserBaseOpsService("go.micro.srv.register", service.Client()),
	})

	handler.Start()

	// Run server
	if err := service.Run(); err != nil {
		fmt.Println("service run",err)
		log.Fatal(err)
	}
}
