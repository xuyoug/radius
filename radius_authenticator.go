package radius

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
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

//SetAuthenticator设置radius的authenticator
func (r *Radius) SetAuthenticator(secret string) {
	//
	if r.Code == CodeAccessRequest {
		return
	}

	if r.Code == CodeAccountingRequest || r.IsRespons() || r.IsCOARespons() {
		buf := bytes.NewBuffer([]byte{})
		if r.Code == CodeAccountingRequest {
			r.Authenticator = Authenticator([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		}
		r.Write(buf)
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
func (r *Radius) IsRequestValid(secret string) bool {
	//
	if r.Code == CodeAccountingRequest {
		buf := bytes.NewBuffer([]byte{})
		authtcr := r.Authenticator
		r.Authenticator = Authenticator([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		r.Write(buf)
		buf.Write([]byte(secret))
		m := md5.Sum(buf.Bytes())
		if Authenticator(m) != authtcr {
			return false
		}
	}
	return true
}

//验证
//对于其它报文全部返回true
func (r *Radius) IsResponseValid(src_authtcr Authenticator, secret string) bool {
	//
	if r.IsRespons() {
		buf := bytes.NewBuffer([]byte{})
		authtcr := r.Authenticator
		r.Authenticator = src_authtcr
		r.Write(buf)
		buf.Write([]byte(secret))
		m := md5.Sum(buf.Bytes())
		if Authenticator(m) != authtcr {
			return false
		}
	}
	return true
}
