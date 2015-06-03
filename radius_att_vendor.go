package radius

import (
	"strconv"
)

//厂商属性类型只有两种：IETF和TYPE4，就不封装了，直接使用字符串

//定义Vender字段类型
type VendorId uint32

//关于VendorId的方法和根据字符串获取VendorId的方法
func (v VendorId) IsValidVendor() bool {
	_, ok := list_vendor_id[v]
	return ok
}

//String 返回venderid的字符串
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

//获取vendor的类型字符串
func (v VendorId) VendorTypestring() string {
	s, ok := list_vendor_id[v]
	if ok {
		return s.Type
	}
	return ""
}

//对外提供由字符串获取vendorid的方法
func GetVendorId(s string) (VendorId, error) {
	s = stringfix(s)
	v, ok := list_vendor_name[s]
	if ok {
		return v, nil
	}
	return VENDOR_NO, ERR_VENDOR_INVALNAME
}

//该vendorid下依据属性名称获取完整属性表达
func (v VendorId) getAttByName(s string) (AttIdV, error) {
	s = stringfix(s)
	vtype := v.VendorTypestring()
	if vtype == "" || !v.IsValidVendor() {
		return ATTIDV_ERR, ERR_ATT_UNK
	}
	if vtype == "TYPE4" {
		vid, ok := list_AttV4_name[v][s]
		if ok {
			return AttIdV{v, vid}, nil
		}
		return ATTIDV_ERR, ERR_ATT_UNK
	}
	if vtype == "IETF" {
		vid, ok := list_AttV_name[v][s]
		if ok {
			return AttIdV{v, vid}, nil
		}
		return ATTIDV_ERR, ERR_ATT_UNK
	}
	return ATTIDV_ERR, ERR_ATT_UNK
}

//该vendorid下依据属性序号获取完整属性表达
func (v VendorId) getAttById(i int) (AttIdV, error) {
	vtype := v.VendorTypestring()
	if vtype == "" || !v.IsValidVendor() {
		return ATTIDV_ERR, ERR_ATT_UNK
	}
	if vtype == "TYPE4" {
		vid := AttV4(i)
		_, ok := list_AttV4_id[v][vid]
		if ok {
			return AttIdV{v, vid}, nil
		}
		return ATTIDV_ERR, ERR_ATT_UNK
	}
	if vtype == "IETF" {
		vid := AttV(i)
		_, ok := list_AttV_id[v][vid]
		if ok {
			return AttIdV{v, vid}, nil
		}
		return ATTIDV_ERR, ERR_ATT_UNK
	}
	return ATTIDV_ERR, ERR_ATT_UNK
}
