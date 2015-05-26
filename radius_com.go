package radius

import (
	"fmt"
	"io"
	"strings"
)

func strpredone(s string) string {
	ss := strings.ToUpper(s)
	ss = strings.TrimSpace(ss)
	return ss
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("Fatal error: %s", err.Error())
	}
}

func isEOF(err error) bool {
	if err == io.EOF {
		return true
	} else {
		return false
	}
}
