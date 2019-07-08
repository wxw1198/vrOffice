package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/go-ini/ini"
	log4plus "github.com/alecthomas/log4go"
)

var Log log4plus.Logger
var ConfigFile *ini.File

func init() {
	var err error
	ConfigFile,err = ini.Load("config.ini","")
	Log = log4plus.NewLogger()

	err = os.MkdirAll("/var/log/go", os.ModePerm)
	if err != nil {
		fmt.Println("dealCompress MkdirAll, err:%v", err)
	}

	flw := log4plus.NewFileLogWriter("/var/log/go/"+"vroffice", true)
	flw.SetRotateDaily(true)
	Log.AddFilter("file", log4plus.DEBUG, flw) //输出到文件,级别为DEBUG,文件名为test.log,每次追加该原文件
}

//获取一个随机数
func GetRandNum() string {
	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(randNum.Intn(1000000))
}

//清零空的字符
func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

//清理空格 回车 换行
func ClearSpecialChar(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)

	return str
}

// 返回true，路径存在
// 返回false并且无错，路径不存在
// 返回错误，不确定路径是否存在
func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DealPanic() {
	var err error
	r := recover()
	if r != nil {
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("Unknown error")
		}

		Log.Error("in go process panic: %s", err.Error())
		fmt.Println("in go process panic:", err.Error())

		Log.Error("%s", debug.Stack())
	}
}

// 获得目录路径。若路径为空，错误；若不含/，返回当前路径'.'
// 例：若path为"/a/b/c/d/index.m3u8////"，dirName为"/a/b/c/d"
func GetDirName(path string) (dirName string, err error) {
	i := len(path) - 1
	if i < 0 {
		return "", errors.New("path is nil")
	}
	// Remove trailing slashes
	for ; i > 0 && path[i] == '/'; i-- {
		path = path[:i]
	}

	i = strings.LastIndex(path, "/")

	if i < 0 {
		return ".", nil
	}

	dirName = path[:i]
	return dirName, nil
}

// 对外暴露的接口：
// 将http请求中的包体转换为json格式，但是不清空包体内容，还可以重新使用包体
func ParseReqBodyToJsonUnclosed(r *http.Request, bodyStruct interface{}, delBody bool) bool {

	if r.ContentLength <= 0 {
		return false
	}
	var bodySlc []byte = make([]byte, 1024)
	bodyLen, readErr := r.Body.Read(bodySlc)
	bodySlc = bodySlc[:bodyLen]
	str := string(bodySlc)
	Log.Info("str =%s", str)

	if readErr == nil {
		Log.Info("read body error:", readErr.Error())
	} else {
		Log.Info("the body has ", bodyLen, " bytes")
	}
	err := json.Unmarshal([]byte(str), bodyStruct)
	if err == nil {
		return true
	} else {
		return false
	}
}

//把检索到数据表的内容转化成为JSON格式
func RowsToJson(rows *sql.Rows) string {
	columns, err := rows.Columns()
	if err != nil {
		Log.Error("row to json:%s", err)
		return ""
	}

	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	var jsonData []byte
	if jsonData, err = json.Marshal(tableData); err != nil {
		fmt.Println("2 row to json:", err)
		return ""
	}

	//fmt.Println(string(jsonData))
	return string(jsonData)
}

//echo -n "value" | openssl sha1 -hmac "key"
//hmac_sha1 此函数名是对应外在的函数名,作用类似md5
func HmacSha1(content, key string) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(content))

	return mac.Sum(nil)
}
