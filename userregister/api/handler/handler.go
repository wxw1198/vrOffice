package handler

import (
	"container/list"
	"context"
	"github.com/micro/go-micro/errors"
	"log"
	"sync"

	proto "github.com/wxw1198/vrOffice/userregister/proto"
)

var (
	registerId  = "100" //待修改
	errRegistering = errors.New(registerId, "the account is registering", 0)
)

type taskInfo struct {
	taskList  *list.List
	taskFound map[string]struct{} //手机号作为Key
	mutex     sync.Mutex
}

var ti taskInfo

func Start() {
	ti = taskInfo{taskList: new(list.List)}
}

type RegisterHandler struct {
	//Client proto.RegisterService
}

//请求入队列，然后此次请求处于等待状态，直到此请求被处理结束后，才返回
//注册信息：手机号，邮箱，昵称，密码
func (g *RegisterHandler) RegisterUser(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	log.Print("Received Greeter.Hello API request", req.Name)

	ti.mutex.Lock()
	defer ti.mutex.Unlock()
	//0 检查是否请求已经在队列里面
	if _, ok := ti.taskFound[req.MobileNum]; ok {
		return errRegistering
	}

	// 1 检查是否已经注册
	// 2 检查注册参数
	// 3 入队列，然后等待
	ti.taskList.PushBack(req)
	ti.taskFound[req.MobileNum] = struct{}{}

	reqChan := make(chan bool)
	go func() {
		//3 开协程，启动信息入库

		reqChan <- true
	}()

	// 4 等待注册结束
	<- reqChan

	// set api response
	rsp.Msg = "register success"
	return nil
}
