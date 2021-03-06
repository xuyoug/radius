package radius

import (
	"errors"
)

var (
	//格式化错误
	ERR_RADIUS_FMT = errors.New("Format radius package error")

	//系统错误，不支持的radius类型
	ERR_NOTSUPPORT = errors.New("Not a support radius type")
	//radius长度错误
	ERR_LEN_INVALID = errors.New("Invalid radius Length")

	//vendor错误
	ERR_VENDOR_INVALID   = errors.New("Invalid radius Vendor id")
	ERR_VENDOR_INVALNAME = errors.New("Invalid radius Vendor name")
	ERR_VENDOR_SET       = errors.New("Set Vendor error")

	//属性错误
	ERR_ATT_FMT   = errors.New("Attribute format error")
	ERR_ATT_UKN   = errors.New("Unknow Attribute")
	ERR_ATT_TYPE  = errors.New("Error Attribute type")
	ERR_ATT_VTYPE = errors.New("Error Attribute value type")
	ERR_ATT_SET   = errors.New("Error on set Attribute value")
	ERR_ATT_NO    = errors.New("No such Attribute")

	//其它
	ERR_VALUE_TYPE = errors.New("Invalid value type")

	//具体属性错误
	ERR_USERNAME_INVALID = errors.New("Invalid username")
	ERR_USERNAME_NUL     = errors.New("Username is null")
	ERR_PASSWD_NUL       = errors.New("Password is null")

	//设置属性错误
	ERR_SET_ATTR = errors.New("Set Attribute error")
)
