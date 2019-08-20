package db

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/wxw1198/vrOffice/userregister/proto"
	"log"
	"time"
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
		fmt.Printf("get score failed, err:%v\n", err)
		return false
	}

	req := proto.Request{}
	json.Unmarshal([]byte(val),&req)

	return req.MobileNum == mobileNum
}


func (u userRegisterRedis) RegisterToDB(req *proto.Request) {
	val, _ := json.Marshal(req)

	u.redisdb.Set(req.Name, val, time.Minute*12)
	u.redisdb.Set(req.MobileNum, val, time.Minute*12)
}
