package radiussvr

//密钥管理结构体定义
type SecretList struct {
	default_secret string
	list_secret    map[string]string
}

//根据默认密钥初始化密钥管理器
func NewSecretList(defaultsecret string) *SecretList {
	sl := new(SecretList)
	sl.default_secret = defaultsecret
	sl.list_secret = make(map[string]string)
	return sl
}

//设置ip及其对应密钥
func (sl *SecretList) SetSecret(ip string, secret string) {
	sl.list_secret[ip] = secret
}

//判断是否有该ip的记录
func (sl *SecretList) IsValidNode(ip string) bool {
	_, ok := sl.list_secret[ip]
	return ok
}

//根据IP查询密钥
//无记录则返回默认密钥
func (sl *SecretList) GetSecret(ip string) string {
	s, ok := sl.list_secret[ip]
	if ok {
		return s
	}
	return sl.default_secret
}
