package server

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	"time"

	//hello "github.com/micro/examples/greeter/srv/proto/hello"
	proto "github.com/wxw1198/vrOffice/userregister/proto"
	"github.com/micro/go-micro"

	"context"
)

type RegisterServer struct{}

func (r *RegisterServer) RegisterUser(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Print("Received Say.Hello request")
	rsp.Msg = "Hello ____" + req.Name
	log.Print("recv :", req.Name)

	return nil
}

func main() {
	//reg := etcdv3.NewRegistry(func(op *registry.Options){
	//	op.Addrs = []string{
	//		"http://192.168.3.34:2379", "http://192.168.3.18:2379", "http://192.168.3.110:2379",
	//	}
	//})

	reg := etcdv3.NewRegistry(func(op *registry.Options){
			op.Addrs = []string{
				"47.88.230.122:55556",
			}
		})


	//etcdv3.NewRegistry()
	//mdns.NewMDNSService()
	//zookeeper.NewRegistry()
	//kubernetes.NewRegistry()
	service := micro.NewService(
		micro.Name("go.micro.srv.greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(reg),
		//micro.Address("127.0.0.1:5555"),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	proto.RegisterRegisterHandler(service.Server(), new(RegisterServer))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
