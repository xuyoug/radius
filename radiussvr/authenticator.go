package radiussvr

//计算计费请求报文的Authenticator有效性
//记帐请求包中的鉴别码中包含对一个由编码＋标识符＋长度＋16个为0的八位字节＋请求属性＋共享密钥所构成的八位字节流进行MD5哈希计算得到的代码
func (sr *SrcRadius) CheckAuthenticator() bool {
	return sr.Radius.IsRequestValid(sr.Secret)
}
