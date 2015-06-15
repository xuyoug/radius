package radius

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

//stringfix预处理字符串
//将字符串转换为大写
//为了效率，不去除空格
func stringfix(s string) string {
	ss := strings.ToUpper(s)
	return ss
}

//checkErr检查错误，如有错误则打印
//用于调试
func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error: %s", err.Error())
	}
}

//isEOF检查错误是否是io.EOF
func isEOF(err error) bool {
	if err == io.EOF {
		return true
	} else {
		return false
	}
}

//随机种子
var rand_source rand.Source
var r_rand *rand.Rand

//初始化随机种子
func init() {
	rand_source = rand.NewSource(int64(time.Now().Nanosecond()))
	r_rand = rand.New(rand_source)
}

//计算随机数
func RandInt(i int) int {
	return r_rand.Intn(i)
}

//计算随机数
func RandBit() byte {
	return byte(r_rand.Intn(256))
}
