package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/wxw1198/vrOffice/userregister/proto"

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

	if iniFilePath == ""{
		iniFilePath = "config.ini"
	}

	conFile, err := ini.Load(iniFilePath, "")
	sec, err := conFile.GetSection("mysql")
	if err != nil {
		utils.Log.Error("ini file ,mysql is not set")
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
		utils.Log.Error(err.Error())
		return false, nil
	}
	return true, db
}

func (t usersTbl) Insert(username, password, server string) bool {
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	fmt.Println("server:", server)

	autu := utils.HmacSha1(username, password)
	encodedSign := base64.RawURLEncoding.EncodeToString(autu)

	insertSql := fmt.Sprintf("INSERT INTO %s(username,password,authorization, servername) VALUES('%s','%s','%s','%s')", t.tableName, username, password, encodedSign, server)
	fmt.Println("insertSql:", insertSql)

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
		utils.Log.Error("dbExec:%s", err.Error())
		return false
	}

	nums, _ := ret.RowsAffected()
	id, _ := ret.LastInsertId()

	fmt.Println("dbExec:id", id, "num:", nums)

	return true
}

//返回map ,key 是authorization, value 是 servername
func (t usersTbl) SelectAuthAndServerName() map[string]string {
	var sqlQuery string

	sqlQuery = fmt.Sprintf("SELECT authorization,servername FROM %s", t.tableName)

	fmt.Println(sqlQuery)

	return dbQueryAuthAndServername(t.dbConnectParam, sqlQuery)
}

//处理查询命令，结果返回json格式数据
func dbQueryAuthAndServername(connectParam, sql string) map[string]string {
	retSets := make(map[string]string)

	//如果有多个defer表达式，调用顺序类似于栈，越后面的defer表达式越先被调用。
	b, db := openDB(connectParam)
	if b == false {
		return retSets
	}
	defer db.Close()

	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("selectAllContent err:%v", err)
		return retSets
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("row to json:", err)
		return retSets
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		v0, ok := values[0].([]byte)
		if ok {
			if v1, ok := values[1].([]byte); ok {
				retSets[string(v0)] = string(v1)
				continue
			}
		}

		utils.Log.Error("select table err:%s", connectParam)
		return retSets
	}

	return retSets
}

//返回map ,key 是authorization, value 是 servername
func (t usersTbl) SelectAll() string {
	var sqlQuery string

	sqlQuery = fmt.Sprintf("SELECT username,authorization,servername FROM %s", t.tableName)

	fmt.Println(sqlQuery)

	return dbQuery(t.dbConnectParam, sqlQuery)
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
		fmt.Println("selectAllContent err:%v", err)
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


func (t usersTbl)MobileNumExist(mobileNum string)bool{
	return false
}


func (t usersTbl)UserNameExist(name string)bool {

	return false
}

func (t usersTbl)RegisterToDB(req *proto.Request)  {
	// todo
}