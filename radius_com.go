package radius

import (
	"fmt"
	"io"
	"strings"
)

//stringfix预处理字符串
//将字符串转换为大写，并去除空格
func stringfix(s string) string {
	ss := strings.ToUpper(s)
	ss = strings.TrimSpace(ss)
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
