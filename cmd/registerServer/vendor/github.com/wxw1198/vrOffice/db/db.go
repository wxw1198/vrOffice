package db

import (
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
)

type dbServer struct {
	diskDBTokenTbl *tokenTbl
	diskDBUserTbl  *usersTbl
	cache          *userRegisterRedis
}

func NewDBServer() *dbServer {
	return &dbServer{
		diskDBUserTbl: NewTable(""),
		cache:         NewUserRegister(),
	}
}

func (db dbServer) MobileNumExist(mobileNum string) bool {
	b := db.cache.MobileNumExist(mobileNum)
	if !b {
		return db.diskDBUserTbl.MobileNumExist(mobileNum)
	}

	return true
}

//注册时，要求强一致性
func (db dbServer) RegisterToDB(req *proto.RegRequest) bool {
	db.diskDBUserTbl.RegisterToDB(req)
	db.cache.RegisterToDB(req)

	return true
}

//注销时，要求强一致性
func (db dbServer) UnRegisterFromDB(req *proto.UnRegRequest) bool {

	db.diskDBUserTbl.UnRegisterFromDB(req)
	db.cache.UnRegisterFromDB(req)

	return true
}

func (db dbServer) CheckUserInfo(request *proto.LoginRequest) bool {
	if db.cache.CheckUserInfo(request) {
		return true
	}

	if b, regReq := db.diskDBUserTbl.CheckUserInfo(request); b {
		go db.cache.SyncRegReqData(regReq)
		return true
	}

	return false
}

func (db dbServer) StoreLoginToken(mobileNum, token string) {

	db.diskDBTokenTbl.StoreLoginToken(mobileNum, token)
	db.cache.StoreLoginToken(mobileNum, token)
}

func (db dbServer) ExistToken(token string) bool {
	if db.cache.ExistToken(token) {
		return true
	}

	if db.diskDBTokenTbl.ExistToken(token) {
		return true
	}

	return false
}

func (db dbServer) DelLoginToken(token string) {
	db.cache.DelLoginToken(token)

	db.diskDBTokenTbl.DelLoginToken(token)
}
