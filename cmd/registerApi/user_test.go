package main

import (
	"encoding/json"
	"fmt"
	"github.com/wxw1198/vrOffice/userbaseoperation/proto"
	"testing"
)

func TestRegRequest_XXX_Marshal(t *testing.T) {
	rr := proto.RegRequest{MobileNum: "18019247756", Name:"12345678",Password:"123456"}
	b,_ := json.Marshal(rr)
	fmt.Println(string(b))
}
