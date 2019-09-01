package handler

import (
	"context"
	"encoding/json"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
	"log"
)

type logoutTaskInfo struct {
	taskList  chan *proto.LogoutRequest
	taskFound map[string]struct{} //手机号作为Key
	resultCh  chan struct{}
}

var loti logoutTaskInfo

//请求入队列，然后此次请求处于等待状态，直到此请求被处理结束后，才返回
//注册信息：手机号，邮箱，昵称，密码
//业务逻辑放到此模块处理
func (g *UserBaseOperationHandler) Logout(ctx context.Context, req *proto.LogoutRequest, rsp *proto.LogoutResponse) error {
	log.Print("Received Greeter.Hello API request", req.Token)

	// 0 检查注册参数
	if req.Token == "" {
		return errParam
	}

	//1 检查是否请求已经在队列里面
	if _, ok := ti.taskFound[req.Token]; ok {
		return errRegistering
	}

	// 3 入队列，然后等待
	loti.taskList <- req
	loti.taskFound[req.Token] = struct{}{}

	reqChan := make(chan bool)
	go func() {
		select {
		case e := <-loti.taskList:
			response, err := g.Client.Logout(ctx, e)
			if err != nil {
				log.Fatal(err)
				rsp.Msg = err.Error()
			}

			b, _ := json.Marshal(map[string]string{
				"message": response.Msg,
			})
			rsp.Msg = string(b)
			//3 开协程，启动信息入库
			reqChan <- true
			return
		}
	}()

	// 4 等待注册结束
	<-reqChan
	return nil
}
