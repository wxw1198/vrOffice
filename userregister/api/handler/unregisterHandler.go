package handler

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/micro/go-micro/errors"
	"github.com/wxw1198/vrOffice/userregister/proto"
)

var (
	errDeleting = errors.New(registerId, " unregister error", 4)
)

type dbUnRegisterInterface interface {
	MobileNumExist(string) bool
	UserNameExist(string) bool
	UnRegisterFromDB(request *proto.UnRegRequest)
}


var dti delTaskInfo

func (g *RegisterHandler) UnRegisterUser(ctx context.Context, req *proto.UnRegRequest, rsp *proto.UnRegResponse) error {
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
	if _, ok := dti.taskFound[req.MobileNum]; ok {
		return errDeleting
	}

	// 3 入队列，然后等待
	dti.taskList <- req
	dti.taskFound[req.MobileNum] = struct{}{}

	reqChan := make(chan bool)
	go func() {
		select {
		case e := <-dti.taskList:
			response, err := g.Client.UnRegisterUser(ctx, e)
			if err != nil {
				log.Fatal(err)
				rsp.Msg = err.Error()
			}

			b, _ := json.Marshal(map[string]string{
				"message": response.Msg,
			})

			rsp.Msg = string(b)
			reqChan <- true
			return
		}
	}()

	// 4 等待注册结束
	<-reqChan

	return nil
}
