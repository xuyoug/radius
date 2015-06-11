package radiussvr

import (
	"bytes"
	"crypto/md5"
	"github.com/xuyoug/radius"
)

//计算响应报文的Authenticator
//ResponseAuth=MD5(Code+ID+Length+RequestAuth+Attributes+Secret)
func (rr *ReplyRadius) ReplyAuthenticator() []byte {
	rr.buf.Write([]byte(rr.Secret))
	m := md5.Sum(rr.buf.Bytes())
	return m[:]
}

//计算计费请求报文的Authenticator有效性
//记帐请求包中的鉴别码中包含对一个由编码＋标识符＋长度＋16个为0的八位字节＋请求属性＋共享密钥所构成的八位字节流进行MD5哈希计算得到的代码
func (sr *SrcRadius) IsValidAuthenticator() bool {
	if sr.Radius.R_Code == radius.CodeAccountingRequest {
		tmp_a := sr.Radius.R_Authenticator
		sr.Radius.R_Authenticator = radius.R_Authenticator([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		tmp_buf := bytes.NewBuffer([]byte{})
		sr.Radius.WriteToBuff(tmp_buf)
		tmp_buf.Write([]byte(sr.Secret))
		m := md5.Sum(tmp_buf.Bytes())
		b := []byte(tmp_a)
		for i := 0; i < 16; i++ {
			if m[i] != b[i] {
				return false
			}
		}
	}
	return true
}
