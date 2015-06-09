package radiuscli

import (
	"github.com/xuyoug/radius"
	"sync"
	//"net"
	"math/rand"
	"time"
)

//radius客户端的实现封装

//
func NewAuthAuthenticator() radius.R_Authenticator {
	bs := make([]byte, 16)
	for i := 0; i < 16; i++ {
		bs = append(bs, byte(getrand(255)))
	}
	return radius.R_Authenticator(bs)
}

//计算随机数
func getrand(i int) int {
	return cli_rand.Intn(i)
}

var cli_source rand.Source
var cli_rand *rand.Rand

func init() {
	cli_source = rand.NewSource(int64(time.Now().Nanosecond()))
	cli_rand = rand.New(cli_source)
}
