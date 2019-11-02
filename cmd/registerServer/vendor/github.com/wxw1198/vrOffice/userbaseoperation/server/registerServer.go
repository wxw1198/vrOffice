package server

//数据服务器，负责数据写入到REDIS MYSQL

import (
	"context"
	"crypto/sha256"
	"strconv"

	"github.com/micro/go-micro/errors"
	"github.com/wxw1198/vrOffice/db"
	"github.com/wxw1198/vrOffice/log"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
	"github.com/wxw1198/vrOffice/utils"
)

var (
	registerId     = "101" //待修改
	errMobileExist = errors.New(registerId, "mobile registed again", 2)
	errUserExist   = errors.New(registerId, "user is exist", 3)
	errUserInfo    = errors.New(registerId, "login && user info err", 4)
)

type dbRegisterInterface interface {
	MobileNumExist(string) bool
	RegisterToDB(*proto.RegRequest) bool
	UnRegisterFromDB(request *proto.UnRegRequest) bool
	CheckUserInfo(request *proto.LoginRequest) bool
	StoreLoginToken(mobileNum, token string)
	ExistToken(token string) bool
	DelLoginToken(token string)
}

type UserBaseOpsServer struct {
	db dbRegisterInterface
}

func NewUserBaseOpsServer() *UserBaseOpsServer {
	return &UserBaseOpsServer{
		db: db.NewDBServer(),
	}
}

func (r *UserBaseOpsServer) RegisterUser(ctx context.Context, req *proto.RegRequest, rsp *proto.RegResponse) error {

	rsp.Msg = "Hello ____" + req.Name
	log.Debug("recv :", req.Name)

	// 1 检查是否已经注册
	// 1.1 手机号
	if r.db.MobileNumExist(req.MobileNum) {
		rsp.Msg = errMobileExist.Error()
		return errMobileExist
	}

	// 2 数据入库
	b := r.db.RegisterToDB(req)
	if !b {
		log.Error("RegisterToDB fail")
		rsp.Msg = "data to db err"
	}

	return nil
}

func (r *UserBaseOpsServer) UnRegisterUser(ctx context.Context, req *proto.UnRegRequest, rsp *proto.UnRegResponse) error {

	rsp.Msg = "Hello ____" + req.Name
	log.Debug("UnRegisterUser :", req.Name)

	// 1 检查是否已经注册
	// 1.1 手机号
	if r.db.MobileNumExist(req.MobileNum) {
		rsp.Msg = errMobileExist.Error()
		return errMobileExist
	}

	// 2 数据入库
	b := r.db.UnRegisterFromDB(req)
	if !b {
		log.Error("req:%v, UnRegisterFromDB fail", req)
		rsp.Msg = "data to db err"
	}

	return nil
}

func (r *UserBaseOpsServer) Login(ctx context.Context, req *proto.LoginRequest, rsp *proto.LoginResponse) error {
	//验证登录信息
	if !r.db.CheckUserInfo(req) {
		return errUserInfo
	}

	//生成登录token
	rn := utils.RandNum()
	randStr := strconv.FormatUint(rn, 10)
	sum := sha256.Sum256([]byte(randStr))
	rsp.Token = string(sum[:])

	r.db.StoreLoginToken(req.MobileNum, rsp.Token)
	return nil
}

func (r *UserBaseOpsServer) Logout(ctx context.Context, req *proto.LogoutRequest, rsp *proto.LogoutResponse) error {
	//验证登出信息
	if r.db.ExistToken(req.Token) {
		//删除登录信息
		r.db.DelLoginToken(req.Token)
	}

	return nil
}

//检查token是否存在
func (r *UserBaseOpsServer) CheckToken() bool {

	//todo
	return false
}
