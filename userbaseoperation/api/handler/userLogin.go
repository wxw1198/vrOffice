package handler

import (
	"encoding/json"
	"context"
	"log"
	"strings"

	"github.com/wxw1198/vrOffice/userbaseoperation/proto"

)

type loginTaskInfo struct {
	taskList  chan *proto.LoginRequest
	taskFound map[string]struct{} //手机号作为Key
	resultCh  chan struct{}
}

var liti loginTaskInfo


//请求入队列，然后此次请求处于等待状态，直到此请求被处理结束后，才返回
//注册信息：手机号，邮箱，昵称，密码
//业务逻辑放到此模块处理
func (g *UserBaseOperationHandler) Login(ctx context.Context, req *proto.LoginRequest, rsp *proto.LoginResponse) error {
	log.Print("Received Greeter.Hello API request", req.MobileNum)

	// 0 检查注册参数
	checkParm := func() bool {
		if req.MobileNum == "" || len(req.Password) <= 6 {
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

	// 3 入队列，然后等待
	liti.taskList <- req
	liti.taskFound[req.MobileNum] = struct{}{}

	reqChan := make(chan bool)
	go func() {
		select {
		case e := <-liti.taskList:
			response, err := g.Client.Login(ctx, e)
			if err != nil {
				log.Fatal(err)
				rsp.Token = ""
			}

			b, _ := json.Marshal(map[string]string{
				"message": response.Token,
			})
			rsp.Token = string(b)
			//3 开协程，启动信息入库
			reqChan <- true
			return
		}
	}()

	// 4 等待注册结束
	<-reqChan
	return nil
}
