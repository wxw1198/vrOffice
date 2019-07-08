package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wxw1198/vrOffice/utils"
)

type usersTbl struct {
	tableName string
}

var usersDB *sql.DB
var UsersTbl usersTbl

var userDBConnectParam string

func init() {
	UsersTbl.tableName = "usersTable"

	var (
		listenIp string
		listenPort  string
		userName string
		password string
		dbName string
	)
	sec, err := utils.ConfigFile.GetSection("sqlip")
	if err != nil {
		utils.Log.Error("ini file ,sqlip is not set")
		listenIp = "127.0.0.1"
	} else {
		listenIp = sec.Body()
	}

	sec, err = utils.ConfigFile.GetSection("sqlport")
	if err != nil {
		utils.Log.Error("ini file ,sqlport is not set")
		listenPort = "3306"
	} else {
		listenPort = sec.Body()
	}

	sec, err = utils.ConfigFile.GetSection("username")
	if err != nil {
		utils.Log.Error("ini file ,username is not set")
		userName = ""
	} else {
		userName = sec.Body()
	}

	sec, err = utils.ConfigFile.GetSection("password")
	if err != nil {
		utils.Log.Error("ini file ,password is not set")
		password = ""
	} else {
		password = sec.Body()
	}

	sec, err = utils.ConfigFile.GetSection("dbname")
	if err != nil {
		utils.Log.Error("ini file ,password is not set")
		dbName = ""
	} else {
		dbName = sec.Body()
	}

	userDBConnectParam = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", userName, password, listenIp, listenPort, dbName)
	fmt.Println(userDBConnectParam)
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

	return dbExec(userDBConnectParam, insertSql)
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

	return dbQueryAuthAndServername(userDBConnectParam, sqlQuery)
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

	return dbQuery(userDBConnectParam, sqlQuery)
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
	return dbExec(userDBConnectParam, updateSql)
}

func (t usersTbl) UpdatePassword(username, password string) bool {
	autu := utils.HmacSha1(username, password)
	encodedSign := base64.RawURLEncoding.EncodeToString(autu)
	updateSql := fmt.Sprintf("update %s set authorization ='%s',password ='%s' where username='%s'", t.tableName, encodedSign, password, username)
	return dbExec(userDBConnectParam, updateSql)
}


MobileNumExist