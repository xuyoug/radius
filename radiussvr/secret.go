package radiussvr

//
type SecretList struct {
	default_secret string
	list_secret    map[string]string
}

//
func NewSecretList(defaultsecret string) *SecretList {
	sl := new(SecretList)
	sl.default_secret = defaultsecret
	sl.list_secret = make(map[string]string)
	return sl
}

//
func (sl *SecretList) SetSecret(ip string, secret string) {
	sl.list_secret[ip] = secret
}

//
func (sl *SecretList) IsValidNode(ip string) bool {
	_, ok := sl.list_secret[ip]
	return ok
}

//
func (sl *SecretList) GetSecret(ip string) string {
	s, ok := sl.list_secret[ip]
	if ok {
		return s
	}
	return sl.default_secret
}
