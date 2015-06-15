package radius

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
)

//定义authenticator 16字节
type Authenticator [16]byte

const AuthenticatorLength = 16

//methods of R_Authenticator
func (a Authenticator) String() string {
	return fmt.Sprintf("Authenticator %v", a[0:16])
}

//从buffer填充Authenticator
func (a *Authenticator) read(buf *bytes.Buffer) error {
	err := binary.Read(buf, binary.BigEndian, a)
	if err != nil {
		return err
	}
	return nil
}

//将Authenticator写入buffer
func (a Authenticator) write(buf *bytes.Buffer) {
	buf.Write(a[0:16])
}

//newAuthenticator生成随机Authenticator
func newAuthenticator() Authenticator {
	var bs [16]byte
	for i := 0; i < 16; i++ {
		bs[i] = RandBit()
	}
	return Authenticator(bs)
}

//setAuthenticator设置authenticator
func (r *Radius) setAuthenticator(secret string) {
	//
	if r.Code == CodeAccountingRequest || r.IsRespons() || IsCOARespons() {
		buf := bytes.NewBuffer([]byte{})
		if r.Code == CodeAccountingRequest {
			r.Authenticator = Authenticator([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		}
		r.WriteToBuff(buf)
		buf.Write([]byte(secret))
		m := md5.Sum(buf.Bytes())
		r.Authenticator = Authenticator(m)
	} else {
		r.Authenticator = newAuthenticator()
	}
}

//IsAuthenticatorValid鉴别authenticator是否有效
//对于计费响应报文按照协议进行计算验证
//对于其它报文全部返回true
func (r *Radius) IsAuthenticatorValid(secret string) bool {
	//
	if r.Code == CodeAccountingRequest {
		buf := bytes.NewBuffer([]byte{})
		authtcr := r.Authenticator
		r.Authenticator = radius.Authenticator([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		r.WriteToBuff(buf)
		buf.Write([]byte(secret))
		m := md5.Sum(buf.Bytes())
		if Authenticator(m) != authtcr {
			return false
		}
	}
	return true
}
