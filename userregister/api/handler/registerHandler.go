package handler

import (
	"context"
	"log"
	"strings"

	"github.com/micro/go-micro/errors"
	"github.com/wxw1198/vrOffice/userregister/proto"
)

var (
	registerId     = "100" //待修改
	errRegistering = errors.New(registerId, "the account is registering", 0)
	errParam       = errors.New(registerId, "register param err", 1)
	errMobileExist = errors.New(registerId, "mobile registed again", 2)
	errUserExist   = errors.New(registerId, "user is exist", 3)
	errDeleting    = errors.New(registerId, " unregister error", 4)
)

type taskInfo struct {
	taskList  chan *proto.Request
	taskFound map[string]struct{} //手机号作为Key
	resultCh  chan struct{}
}

type delTaskInfo struct {
	taskList  chan *proto.UnRegisterRequest
	taskFound map[string]struct{} //手机号作为Key
	resultCh  chan struct{}
	//mutex     sync.Mutex
}

var ti taskInfo

func Start() {
	ti = taskInfo{taskList: make(chan *proto.Request, 1024)}
	dti = delTaskInfo{taskList: make(chan *proto.UnRegisterRequest, 1024)}
}

type RegisterHandler struct {
	db dbRegisterInterface
	//Client proto.RegisterService
}

type dbRegisterInterface interface {
	MobileNumExist(string) bool
	UserNameExist(string) bool
	RegisterToDB(*proto.Request)
}

//请求入队列，然后此次请求处于等待状态，直到此请求被处理结束后，才返回
//注册信息：手机号，邮箱，昵称，密码
func (g *RegisterHandler) RegisterUser(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Print("Received Greeter.Hello API request", req.Name)

	// 0 检查注册参数
	checkParm := func() bool {
		if req.MobileNum == "" || req.Name == "" || len(req.Password) <= 6 {
			return false
		}
		if len(req.MobileNum) != 11 && strings.HasPrefix(req.MobileNum, "1") == false {
			return false
		}
		return true
	}
	if checkParm() == false {
		return errParam
	}

	//1 检查是否请求已经在队列里面
	if _, ok := ti.taskFound[req.MobileNum]; ok {
		return errRegistering
	}

	// 2 检查是否已经注册
	// 2.1 手机号
	if g.db.MobileNumExist(req.MobileNum) {
		return errMobileExist
	}
	// 2.2 用户名
	if g.db.UserNameExist(req.Name) {
		return errUserExist
	}

	// 3 入队列，然后等待
	ti.taskList <- req
	ti.taskFound[req.MobileNum] = struct{}{}

	reqChan := make(chan bool)
	go func() {
		select {
		case e := <-ti.taskList:
			//3 开协程，启动信息入库
			g.db.RegisterToDB(e)
			reqChan <- true
			return
		}
	}()

	// 4 等待注册结束
	<-reqChan

	// set api response
	rsp.Msg = "register success"
	return nil
}