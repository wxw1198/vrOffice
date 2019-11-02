package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	//fileName := "micro-srv.log"
	//syncWriter := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:  fileName,
	//	MaxSize:   128, //MB
	//	LocalTime: true,
	//	Compress:  true,
	//})
	//encoder := zap.NewProductionEncoderConfig()
	//encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(zap.DebugLevel))
	//log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	//Log = log.Sugar()
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

//把检索到数据表的内容转化成为JSON格式
func RowsToJson(rows *sql.Rows) string {
	columns, err := rows.Columns()
	if err != nil {
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

func RandNum() uint64 {
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand1.Uint64()
}
