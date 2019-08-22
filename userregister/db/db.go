package db

import (
	"github.com/wxw1198/vrOffice/userregister/proto"
)

type dbServer struct {
	diskDB *usersTbl
	cache  *userRegisterRedis
}

func NewDBServer() *dbServer {
	return &dbServer{
		diskDB: NewTable(""),
		cache:  NewUserRegister(),
	}
}

func (db dbServer) MobileNumExist(mobileNum string) bool {
	b := db.cache.MobileNumExist(mobileNum)
	if !b {
		return db.diskDB.MobileNumExist(mobileNum)
	}

	return true
}

//注册时，要求强一致性
func (db dbServer) RegisterToDB(req *proto.Request) bool {
	db.diskDB.RegisterToDB(req)
	db.cache.RegisterToDB(req)

	return true
}


//注销时，要求强一致性
func (db dbServer)UnRegisterFromDB(req *proto.UnRegRequest)bool{

	db.diskDB.UnRegisterFromDB(req)
	db.cache.UnRegisterFromDB(req)

	return true
}