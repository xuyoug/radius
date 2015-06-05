package radiuscli

import (
	"github.com/xuyoug/radius"
	"sync"
	//"net"
	"math/rand"
	"time"
)

//radius客户端的实现封装

//定义客户端的Id序列
var radiuscli_id radius.R_Id

//
var cli_sync sync.Mutex

//
func GetRadiusId() radius.R_Id {
	cli_sync.Lock()
	if radiuscli_id == radius.R_Id(255) {
		radiuscli_id = 0
	} else {
		radiuscli_id++
	}
	cli_sync.Unlock()
	return radiuscli_id
}

func NewAuthenticator() radius.R_Authenticator {

}

//计算随机数

func getrand(i int) int {

}

var cli_source rand.Source
var cli_rand *rand.Rand

func init() {
	cli_source = rand.NewSource(int64(time.Now().Nanosecond()))
	cli_rand = rand.New(cli_source)
}
