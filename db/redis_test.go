package db

import (
	"testing"
	"time"
)

func TestConnectServer(t *testing.T) {

	con := NewUserRegister()
	if con == nil {
		t.Fatal("new user register err")
	}

	con.redisdb.Set("helloredis", "world-vroffice", time.Hour*1)

	ret := con.redisdb.Get("helloredis")

	if ret.Val() != "world-vroffice" {
		t.Fatal("stroage vale err")
	}
}
