package main

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"

	//proto "github.com/micro/examples/greeter/api/rpc/proto/hello"
	proto "github.com/wxw1198/vrOffice/userregister/proto"
	//hello "github.com/micro/examples/greeter/srv/proto/hello"
	"github.com/micro/go-micro"

	"context"
)

type Register struct {
	Client   proto.RegisterService //hello.SayService
}

func (g *Register) RegisterUser(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Print("Received Greeter.Hello API request", req.Name)

	// make the request
	//response, err := g.Client.Hello(ctx, &hello.Request{Name: req.Name})
	//if err != nil {
	//	return err
	//}

	// 1 检查是否已经注册
	// 2 检查注册参数
	//3 开协程，启动信息入库
	// 4 等待注册结束

	// set api response
	rsp.Msg = "register success"
	return nil
}

func main() {
	// Create service
	reg := etcdv3.NewRegistry(func(op *registry.Options){
		op.Addrs = []string{
			"47.88.230.122:55556",
		}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.api.register"),
	)

	// Init to parse flags
	service.Init()

	proto.RegisterRegisterHandler(service.Server(), &Register{
		// Create Service Client
		Client: proto.NewRegisterService("go.micro.api.register", service.Client()),
	})

	// for handler use

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
