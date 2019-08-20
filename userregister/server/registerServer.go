package server


//数据服务器，负责数据写入到REDIS MYSQL

import (
	"context"
	"log"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/wxw1198/vrOffice/userregister/proto"
	"github.com/wxw1198/vrOffice/userregister/db"
)

var (
	registerId     = "101" //待修改
	errMobileExist = errors.New(registerId, "mobile registed again", 2)
	errUserExist   = errors.New(registerId, "user is exist", 3)
)

type dbRegisterInterface interface {
	MobileNumExist(string) bool
	RegisterToDB(*proto.Request)bool
}

type RegisterServer struct {
	db dbRegisterInterface
}

func NewRegisterServer()*RegisterServer  {
	return &RegisterServer{
		db:db.NewDBServer(),
	}
}

func (r *RegisterServer) RegisterUser(ctx context.Context, req *proto.Request, rsp *proto.Response) error {

	rsp.Msg = "Hello ____" + req.Name
	log.Print("recv :", req.Name)

	// 1 检查是否已经注册
	// 1.1 手机号
	if r.db.MobileNumExist(req.MobileNum) {
		rsp.Msg = errMobileExist.Error()
		return errMobileExist
	}

	// 2 数据入库
	b := r.db.RegisterToDB(req)
	if !b {
		rsp.Msg = "data to db err"
	}

	return nil
}

func (r *RegisterServer) UnRegisterUser(ctx context.Context, req *proto.UnRegisterRequest, rsp *proto.UnRegisterResponse) error {

	rsp.Msg = "Hello ____" + req.Name
	log.Print("UnRegisterUser :", req.Name)

	// 1 检查是否已经注册
	// 1.1 手机号
	if r.db.MobileNumExist(req.MobileNum) {
		rsp.Msg = errMobileExist.Error()
		return errMobileExist
	}

	// 2 数据入库
	b := r.db.UnRegisterFromDB(req)
	if !b {
		rsp.Msg = "data to db err"
	}

	return nil
}

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

	//etcdv3.NewRegistry()
	//mdns.NewMDNSService()
	//zookeeper.NewRegistry()
	//kubernetes.NewRegistry()
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
	proto.RegisterRegisterHandler(service.Server(), new(RegisterServer))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
