package server

//数据服务器，负责数据写入到REDIS MYSQL

import (
	"context"
	"log"

	"github.com/micro/go-micro/errors"
	"github.com/wxw1198/vrOffice/userregister/db"
	"github.com/wxw1198/vrOffice/userregister/proto"
	"github.com/wxw1198/vrOffice/utils"
)

var (
	registerId     = "101" //待修改
	errMobileExist = errors.New(registerId, "mobile registed again", 2)
	errUserExist   = errors.New(registerId, "user is exist", 3)
)

type dbRegisterInterface interface {
	MobileNumExist(string) bool
	RegisterToDB(*proto.Request) bool
	UnRegisterFromDB(request *proto.UnRegRequest) bool
}

type RegisterServer struct {
	db dbRegisterInterface
}

func NewRegisterServer() *RegisterServer {
	return &RegisterServer{
		db: db.NewDBServer(),
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
		utils.Log.With(req).Error("RegisterToDB fail")
		rsp.Msg = "data to db err"
	}

	return nil
}

func (r *RegisterServer) UnRegisterUser(ctx context.Context, req *proto.UnRegRequest, rsp *proto.UnRegResponse) error {

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
		utils.Log.With(req).Error("UnRegisterFromDB fail")
		rsp.Msg = "data to db err"
	}

	return nil
}
