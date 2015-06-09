package radius

import (
	"strconv"
)

//厂商属性类型只有两种：IETF和TYPE4，就不封装了，直接使用字符串
//本文件是关于VendorId的方法和根据字符串获取VendorId的方法

//VendorId定义Vender字段类型
type VendorId uint32

//IsValidVendor方法判断是否是有效Vendor
func (v VendorId) IsValidVendor() bool {
	_, ok := list_vendor_id[v]
	return ok
}

//String返回venderid的格式化字符串
func (v VendorId) String() string {
	if v == VENDOR_NO {
		return "VENDOR_NO"
	}
	s, ok := list_vendor_id[v]
	if ok {
		return s.Name + "(" + strconv.Itoa(int(v)) + ")"
	}
	return "UNKNOWN_VENDOR(" + strconv.Itoa(int(v)) + ")"
}

//VendorTypestring获取vendor的类型字符串
func (v VendorId) VendorTypestring() string {
	s, ok := list_vendor_id[v]
	if ok {
		return s.Type
	}
	return ""
}

//GetVendorId提供由vendor名称获取vendorid的方法
func GetVendorId(s string) (VendorId, error) {
	s = stringfix(s)
	v, ok := list_vendor_name[s]
	if ok {
		return v, nil
	}
	return VENDOR_NO, ERR_VENDOR_INVALNAME
}

//getAttByName方法提供该vendorid下依据属性名称获取属性表达
//若该厂商属性有效（已注册）则返回，否则返回错误
func (v VendorId) GetAttByName(s string) (AttIdV, error) {
	s = stringfix(s)
	vtype := v.VendorTypestring()
	if vtype == "" || !v.IsValidVendor() {
		return ATTIDV_ERR, ERR_ATT_UNK
	}

	id, ok := list_attV_name[v][s]
	if ok {
		return AttIdV{v, id}, nil
	}
	return ATTIDV_ERR, ERR_ATT_UNK
}

//getAttById方法提供该vendorid下依据属性序号获取属性表达
//若该厂商属性有效（已注册）则返回，否则返回错误
func (v VendorId) GetAttById(i int) (AttIdV, error) {
	vtype := v.VendorTypestring()
	if vtype == "" || !v.IsValidVendor() {
		return ATTIDV_ERR, ERR_ATT_UNK
	}
	_, ok := list_attV_id[v][i]
	if ok {
		return AttIdV{v, i}, nil
	}
	return ATTIDV_ERR, ERR_ATT_UNK

}
