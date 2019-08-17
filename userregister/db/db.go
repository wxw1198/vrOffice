package db

import (
	"github.com/wxw1198/vrOffice/userregister/proto"
)

type dbServer struct {
	diskDB *usersTbl
	cache * userRegisterRedis
}


func (db dbServer) MobileNumExist(mobileNum string) bool {

	b := db.cache.MobileNumExist(mobileNum)
	if !b {
		return db.diskDB.MobileNumExist(mobileNum)
	}

	return true
}

func (db dbServer) UserNameExist(name string) bool {
	b := db.cache.UserNameExist(name)
	if !b {
		return db.diskDB.UserNameExist(name)
	}

	return true
}

func (db dbServer) RegisterToDB(req *proto.Request) {
	db.cache.RegisterToDB(req)

	db.diskDB.RegisterToDB(req)
}

