package db

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/wxw1198/vrOffice/log"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wxw1198/vrOffice/utils"
)

type usersTbl struct {
	tableName      string
	dbConnectParam string
}

var usersDB *sql.DB

func (t *usersTbl) readConfig(iniFilePath string) error {
	var (
		listenIp   string
		listenPort string
		userName   string
		password   string
		dbName     string
	)

	if iniFilePath == "" {
		iniFilePath = "config.ini"
	}

	conFile, err := ini.Load(iniFilePath, "")
	sec, err := conFile.GetSection("mysql")
	if err != nil {
		log.Error("ini file ,mysql is not set")
		return err
	} else {
		listenIp = sec.Key("sqlip").String()
		if listenIp == "" {
			listenIp = "127.0.0.1"
		}

		listenPort = sec.Key("sqlport").String()
		if listenPort == "" {
			listenPort = "3306"
		}

		userName = sec.Key("username").String()
		password = sec.Key("password").String()
		t.tableName = sec.Key("usertablename").String()
	}

	t.dbConnectParam = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", userName, password, listenIp, listenPort, dbName)
	fmt.Println(t.dbConnectParam)
	return nil
}

func NewTable(iniFilePath string) *usersTbl {
	t := &usersTbl{}
	err := t.readConfig(iniFilePath)
	if err != nil {
		return nil
	}
	return t
}

func openDB(mysqlParam string) (bool, *sql.DB) {
	db, err := sql.Open("mysql", mysqlParam)
	if err != nil { //连接成功 err一定是nil否则就是报错
		fmt.Println(err.Error()) //仅仅是显示异常
		log.Error("openDB:", err.Error())
		return false, nil
	}
	return true, db
}

func (t usersTbl) Insert(username, password, server string) bool {

	autu := utils.HmacSha1(username, password)
	encodedSign := base64.RawURLEncoding.EncodeToString(autu)

	insertSql := fmt.Sprintf("INSERT INTO %s(username,password,authorization, servername) VALUES('%s','%s','%s','%s')", t.tableName, username, password, encodedSign, server)
	log.Debug("insertSql:", insertSql)

	return dbExec(t.dbConnectParam, insertSql)
}

//处理insert /update /del
func dbExec(connectParam, sql string) bool {

	fmt.Println(connectParam)
	fmt.Println(sql)

	b, db := openDB(connectParam)
	if b == false {
		return b
	}
	defer db.Close()

	ret, err := db.Exec(sql)
	if err != nil {
		fmt.Println("dbExec:", err.Error())
		log.Error("dbExec:%s", err.Error())
		return false
	}

	nums, _ := ret.RowsAffected()
	id, _ := ret.LastInsertId()

	log.Debug("dbExec:id", id, "num:", nums)

	return true
}

//处理查询命令，结果返回json格式数据
func dbQuery(connectParam, sql string) string {
	b, db := openDB(connectParam)
	if b == false {
		return ""
	}
	defer db.Close()

	rows, err := db.Query(sql)
	if err != nil {
		log.Error("selectAllContent err:%v", err)
		return ""
	}
	defer rows.Close()

	return utils.RowsToJson(rows)
}

func (t usersTbl) UpdateServer(username, servername string) bool {
	updateSql := fmt.Sprintf("update %s set servername ='%s' where username='%s'", t.tableName, servername, username)
	return dbExec(t.dbConnectParam, updateSql)
}

func (t usersTbl) UpdatePassword(username, password string) bool {
	autu := utils.HmacSha1(username, password)
	encodedSign := base64.RawURLEncoding.EncodeToString(autu)
	updateSql := fmt.Sprintf("update %s set authorization ='%s',password ='%s' where username='%s'", t.tableName, encodedSign, password, username)
	return dbExec(t.dbConnectParam, updateSql)
}

func (t usersTbl) MobileNumExist(mobileNum string) bool {
	sqlQuery := fmt.Sprintf("SELECT mobile,name,password FROM %s where mobile='%s'", t.tableName, mobileNum)
	ret := dbQuery(t.dbConnectParam, sqlQuery)
	if len(ret) > 0 {
		return true
	}
	return false
}

func (t usersTbl) RegisterToDB(req *proto.RegRequest) bool {
	insertSql := fmt.Sprintf("INSERT INTO %s(mobile,name,password) VALUES('%s','%s','%s')", t.tableName, req.MobileNum, req.Name, req.Password)
	log.Debug("insertSql:", insertSql)

	return dbExec(t.dbConnectParam, insertSql)
	// todo
}

func (t usersTbl) UnRegisterFromDB(req *proto.UnRegRequest) bool {
	delSql := fmt.Sprintf("DELETE FROM '%s' WHERE mobile='%s'", t.tableName, req.MobileNum)

	log.Debug("delSql:", delSql)

	return dbExec(t.dbConnectParam, delSql)
}

func (t usersTbl) CheckUserInfo(request *proto.LoginRequest) (bool, *proto.RegRequest) {
	sqlQuery := fmt.Sprintf("SELECT mobile,name,password FROM %s where mobile='%s'", t.tableName, request.MobileNum)
	ret := dbQuery(t.dbConnectParam, sqlQuery)

	regReq := proto.RegRequest{}
	json.Unmarshal([]byte(ret), &regReq)

	return regReq.MobileNum == request.MobileNum && regReq.Password == request.Password, &regReq
}

type tokenTbl struct {
	tableName      string
	dbConnectParam string
}

func (t tokenTbl) StoreLoginToken(mobileNum, token string) {
	// todo
	insertSql := fmt.Sprintf("INSERT INTO %s(mobile,token) VALUES('%s','%s')", t.tableName, token)
	log.Debug("insertSql:", insertSql)

	dbExec(t.dbConnectParam, insertSql)
}

func (t tokenTbl) ExistToken(token string) bool {
	sqlQuery := fmt.Sprintf("SELECT token FROM %s where token='%s'", t.tableName, token)
	ret := dbQuery(t.dbConnectParam, sqlQuery)

	if len(ret) > 0 {
		return true
	}
	return false
}

func (t tokenTbl) DelLoginToken(token string) {
	delSql := fmt.Sprintf("DELETE FROM '%s' WHERE token='%s'", t.tableName, token)

	log.Debug("insertSql:", delSql)

	dbExec(t.dbConnectParam, delSql)
}
