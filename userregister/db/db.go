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


func (db dbServer) RegisterToDB(req *proto.Request) bool {
	db.cache.RegisterToDB(req)

	return db.diskDB.RegisterToDB(req)
}



func (db dbServer)UnRegisterFromDB(req *proto.UnRegisterRequest){
	db.cache.UnRegisterFromDB(req)

	return db.diskDB.UnRegisterFromDB(req)
}