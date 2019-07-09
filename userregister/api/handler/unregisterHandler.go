package handler

import (
	"context"
	"strings"

	"github.com/wxw1198/vrOffice/userregister/proto"
)

type dbUnRegisterInterface interface {
	MobileNumExist(string) bool
	UserNameExist(string) bool
	UnRegisterFromDB(request *proto.UnRegisterRequest)
}
type DeleteHandler struct {
	db dbUnRegisterInterface
	//Client proto.RegisterService
}

var dti delTaskInfo

func (g *DeleteHandler) DeleteUser(ctx context.Context, req *proto.UnRegisterRequest, rsp *proto.Response) error {
	// 0 检查注册参数
	checkParm := func() bool {
		if req.MobileNum == "" || req.Name == "" {
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
		return errDeleting
	}

	// 2 检查是否已经注册
	// 2.1 手机号
	if !g.db.MobileNumExist(req.MobileNum) {
		return errMobileExist
	}
	// 2.2 用户名
	if !g.db.UserNameExist(req.Name) {
		return errUserExist
	}

	// 3 入队列，然后等待
	dti.taskList <- req
	dti.taskFound[req.MobileNum] = struct{}{}

	reqChan := make(chan bool)
	go func() {
		select {
		case e := <-dti.taskList:
			//3 开协程，启动信息入库
			g.db.UnRegisterFromDB(e)
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
