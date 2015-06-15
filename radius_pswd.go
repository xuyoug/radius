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
func Decipher_pap_passwd(psw_bytes []byte, secret string, authenticator_in Authenticator) string {
	if len(psw_bytes)%16 != 0 || len(secret) == 0 {
		return ""
	}
	l := len(psw_bytes) / 16
	psw_out := make([]byte, 0)
	authenticator := [16]byte(authenticator_in)
	var encry_tmp []byte
	for i := 0; i < l; i++ {
		encry := []byte(secret)
		if i == 0 {
			encry = append(encry, authenticator[:]...)
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
func Encry_pap_passwd(psw string, secret string, authenticator_in Authenticator) []byte {
	authenticator := [16]byte(authenticator_in)
	if len(secret) == 0 {
		return []byte{}
	}

	psw_bytes := append([]byte(psw), make([]byte, 16-len(psw)%16)...)
	psw_encried := make([]byte, 0)
	var encry_tmp []byte
	l := len(psw_bytes) / 16
	for i := 0; i < l; i++ {
		encry := []byte(secret)
		if i == 0 {
			encry = append(encry, authenticator[:]...)
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

//
func (r *Radius) CheckPasswd(pwd, secret string) bool {
	if r.Code != CodeAccessRequest {
		return true
	}

	pswd_p := r.GetAttrValue1(ATTID_USER_PASSWORD)
	pswd_c := r.GetAttrValue1(ATTID_CHAP_PASSWORD)
	if pswd_p == nil && pswd_c != nil {
		var clg string
		if c := r.GetAttrValue1(ATTID_CHAP_CHALLENGE); c != nil {
			clg = c.Value().(string)
		} else {
			clg = string([]byte(r.Authenticator[:]))
		}
		if len(pswd_c.Value().(string)) != 17 {
			return false
		}
		chapid := pswd_c.Value().(string)[0:1]
		chapc := pswd_c.Value().(string)[1:17]

		m := md5.New()
		m.Write([]byte(chapid))
		m.Write([]byte(pwd))
		m.Write([]byte(clg))
		m_out := m.Sum(nil)
		for i := 0; i < 16; i++ {
			if m_out[i] != chapc[i] {
				return false
			}
		}
		return true
	}
	if pswd_c == nil && pswd_p != nil {
		return Decipher_pap_passwd([]byte(pswd_p.Value().(string)), secret, r.Authenticator) == pwd
	}
	return false
}

//
func (r *Radius) AddPwd(pwd, secret string, ispap bool) error {
	if r.Code != CodeAccessRequest {
		return ERR_SET_ATTR
	}
	if ispap {
		r.AddAttr(&Attribute{ATTID_USER_PASSWORD, STRING(Encry_pap_passwd(pwd, secret, r.Authenticator))})
	} else {
		clg := []byte(r.Authenticator[:])
		chapid := byte(RandInt(255))

		m := md5.New()
		m.Write([]byte{chapid})
		m.Write([]byte(pwd))
		m.Write(clg)
		m_out := m.Sum(nil)

		passwd := make([]byte, 17)
		passwd[0] = chapid
		for i := 1; i < 17; i++ {
			passwd[i] = m_out[i-1]
		}
		r.AddAttr(&Attribute{ATTID_CHAP_PASSWORD, STRING(passwd)})
	}
	return nil
}
