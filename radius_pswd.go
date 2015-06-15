package radius

import (
	"bytes"
	"crypto/md5"
)

//bytesxor返回两个字节切片的异或
//按长度最短的截取
func bytesxor(i1, i2 []byte) []byte {
	l := len(i2)
	if len(i1) < l {
		l = len(i1)
	}
	o := make([]byte, 0)
	for i := 0; i < l; i++ {
		o = append(o, i1[i]^i2[i])
	}
	return o
}

//Decipher_pap_passwd解密pap密码
//可对外导出
func Decipher_pap_passwd(psw_bytes []byte, secret string, authenticator []byte) string {
	if len(psw_bytes)%16 != 0 || len(secret) == 0 {
		return ""
	}
	l := len(psw_bytes) / 16
	psw_out := make([]byte, 0)
	var encry_tmp []byte
	for i := 0; i < l; i++ {
		encry := []byte(secret)
		if i == 0 {
			encry = append(encry, authenticator...)
		} else {
			encry = append(encry, encry_tmp...)
		}
		//m := crypto.Hash(crypto.MD5).New()
		m := md5.Sum(encry)
		encry_tmp = psw_bytes[i*16 : (i+1)*16]
		psw_out = append(psw_out, bytesxor(m[:], encry_tmp)...)
	}
	return string(bytes.TrimRight(psw_out, string([]byte{0})))
}

//Encry_pap_passwd加密pap密码
//可对外导出
func Encry_pap_passwd(psw string, secret string, authenticator []byte) []byte {
	if len(authenticator) != 16 || len(secret) == 0 {
		return []byte{}
	}
	psw_bytes := append([]byte(psw), make([]byte, 16-len(psw)%16)...)
	psw_encried := make([]byte, 0)
	var encry_tmp []byte
	l := len(psw_bytes) / 16
	for i := 0; i < l; i++ {
		encry := []byte(secret)
		if i == 0 {
			encry = append(encry, authenticator...)
		} else {
			encry = append(encry, encry_tmp...)
		}
		//m := crypto.Hash(crypto.MD5).New()
		m := md5.New()
		m.Write(encry)
		encry_tmp = bytesxor(m.Sum(nil), psw_bytes[i*16:(i+1)*16])
		psw_encried = append(psw_encried, encry_tmp...)
	}
	return psw_encried
}
