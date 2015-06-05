package radiussvr

import (
	"net"
)

//
var SecretList struct {
	default_secret string
	list_secret    map[net.IP]string
}

//
func NewSecretList(defaultsecret string) *SecretList {
	sl := new(SecretList)
	sl.default_secret = defaultsecret
	sl.list_secret = make(map[net.IP]string)
	return sl
}

//
func (sl *SecretList) SetSecret(ip net.IP, secret string) {
	sl.list_secret[ip] = secret
}

//
func (sl *SecretList) IsValidNode(ip net.IP) bool {
	_, ok := sl.list_secret[ip]
	return ok
}

//
func (sl *SecretList) GetSecret(ip net.IP) string {
	s, ok := sl.list_secret[ip]
	if ok {
		return s
	}
	return sl.default_secret
}
