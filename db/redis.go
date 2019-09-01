package db

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/wxw1198/vrOffice/log"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
)

type userRegisterRedis struct {
	// 声明一个全局的redisdb变量
	redisdb *redis.Client
}

func NewUserRegister() *userRegisterRedis {
	u := &userRegisterRedis{
		redisdb: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}

	_, err := u.redisdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return u
}

func (u userRegisterRedis) MobileNumExist(mobileNum string) bool {
	val, err := u.redisdb.Get(mobileNum).Result()
	if err != nil {
		log.Error("get score failed, err:%v\n", err)
		return false
	}

	req := proto.RegRequest{}
	json.Unmarshal([]byte(val), &req)

	return req.MobileNum == mobileNum
}

func (u *userRegisterRedis) RegisterToDB(req *proto.RegRequest) {
	val, _ := json.Marshal(req)

	//u.redisdb.Set(req.Name, val, time.Minute*12)
	u.redisdb.Set(req.MobileNum, val, time.Minute*120)
}

func (u userRegisterRedis) UnRegisterFromDB(req *proto.UnRegRequest) {
	u.redisdb.Del(req.MobileNum)
}

func (u userRegisterRedis) CheckUserInfo(request *proto.LoginRequest) bool {
	val, err := u.redisdb.Get(request.MobileNum).Result()

	if err != nil {
		log.Error("get score failed, err:%v\n", err)
		return false
	}

	req := proto.RegRequest{}
	json.Unmarshal([]byte(val), &req)

	return request.Password == req.Password && request.MobileNum == req.MobileNum
}

//mysql的数据同步到cache
func (u *userRegisterRedis) SyncRegReqData(req *proto.RegRequest) {
	u.RegisterToDB(req)
}

func (u *userRegisterRedis) StoreLoginToken(mobileNum, token string) {
	//u.redisdb.Set(req.Name, val, time.Minute*12)
	u.redisdb.Set(token, mobileNum, time.Minute*120)
}

func (u userRegisterRedis) ExistToken(token string) bool {
	if u.redisdb.Get(token) == nil {
		return false
	} else {
		return true
	}
}

func (u userRegisterRedis)DelLoginToken(token string){
	u.redisdb.Del(token)
}