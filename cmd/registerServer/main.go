package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
	"github.com/wxw1198/vrOffice/userbaseoperation/server"
	"github.com/wxw1198/vrOffice/utils"

	"time"
)

func main() {
	//reg := etcdv3.NewRegistry(func(op *registry.Options){
	//	op.Addrs = []string{
	//		"http://192.168.3.34:2379", "http://192.168.3.18:2379", "http://192.168.3.110:2379",
	//	}
	//})

	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			"47.88.230.122:55556",
		}
	})

	service := micro.NewService(
		micro.Name("go.micro.srv.register"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(reg),
		//micro.Address("127.0.0.1:5555"),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	//proto.RegisterRegisterHandler(service.Server(), new(UserBaseOpsServer))

	proto.RegisterUserBaseOpsHandler(service.Server(), server.NewUserBaseOpsServer())

	// Run server
	if err := service.Run(); err != nil {
		utils.Log.Fatal(err)
	}
}
