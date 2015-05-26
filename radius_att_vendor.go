package radius

import (
	"strconv"
)

//厂商属性类型只有两种：IETF和TYPE4，就不封装了，直接使用字符串

//定义Vender字段类型
type VendorId uint32

func (v VendorId) ToInt() int {
	return int(v)
}

//定义厂商列表
const (
	VENDOR_NO         VendorId = 0 //本包中  此vendorid作为标准属性的封装，故NAME，STRING、ISVAILD方法中不认为是错误
	VENDOR_ACC        VendorId = 5
	VENDOR_CISCO      VendorId = 9
	VENDOR_XYLOGICS   VendorId = 15
	VENDOR_MERIT      VendorId = 61
	VENDOR_GANDALF    VendorId = 64
	VENDOR_SHIVA      VendorId = 166
	VENDOR_LIVINGSTON VendorId = 307
	VENDOR_MICROSOFT  VendorId = 311
	VENDOR_3COM       VendorId = 429
	VENDOR_ASCEND     VendorId = 529
	VENDOR_BAY        VendorId = 1584
	VENDOR_LUCENT     VendorId = 1751
	VENDOR_REDBACK    VendorId = 2352
	VENDOR_APTIS      VendorId = 2634
	VENDOR_MASTERSOFT VendorId = 5401
	VENDOR_QUINTUM    VendorId = 6618
	VENDOR_HUAWEI     VendorId = 2011
	VENDOR_JUNIPER    VendorId = 4874
	VENDOR_ZTE        VendorId = 3902
	VENDOR_ALCATEL    VendorId = 6527
)

//
type const_vendor struct {
	Name string
	Type string
}

//定义VEBDOR的ID映射
var list_vendor_id map[VendorId]const_vendor = map[VendorId]const_vendor{
	5:    {"ACC", "IETF"},
	9:    {"CISCO", "IETF"},
	15:   {"XYLOGICS", "IETF"},
	61:   {"MERIT", "IETF"},
	64:   {"GANDALF", "IETF"},
	166:  {"SHIVA", "IETF"},
	307:  {"LIVINGSTON", "IETF"},
	311:  {"MICROSOFT", "IETF"},
	429:  {"3COM", "TYPE4"},
	529:  {"ASCEND", "IETF"},
	1584: {"BAY", "IETF"},
	1751: {"LUCENT", "IETF"},
	2352: {"REDBACK", "IETF"},
	2634: {"APTIS", "IETF"},
	5401: {"MASTERSOFT", "IETF"},
	6618: {"QUINTUM", "IETF"},
	2011: {"HUAWEI", "IETF"},
	4874: {"JUNIPER", "IETF"},
	3902: {"ZTE", "IETF"},
	6527: {"ALCATEL", "IETF"},
}

//定义VENDOR的NAME映射
var list_vendor_name map[string]VendorId = map[string]VendorId{
	"ACC":        5,
	"CISCO":      9,
	"XYLOGICS":   15,
	"MERIT":      61,
	"GANDALF":    64,
	"SHIVA":      166,
	"LIVINGSTON": 307,
	"MICROSOFT":  311,
	"3COM":       429,
	"ASCEND":     529,
	"BAY":        1584,
	"LUCENT":     1751,
	"REDBACK":    2352,
	"APTIS":      2634,
	"MASTERSOFT": 5401,
	"QUINTUM":    6618,
	"HUAWEI":     2011,
	"JUNIPER":    4874,
	"ZTE":        3902,
	"ALCATEL":    6527,
}

//关于VendorId的方法和根据字符串获取VendorId的方法
func (v VendorId) IsvalidVendor() bool {
	if v == VENDOR_NO {
		return true
	}
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
func (v VendorId) Typestring() string {
	s, ok := list_vendor_id[v]
	if ok {
		return s.Type
	}
	return ""
}

//对外提供由字符串获取vendorid的方法
func GetVendorId(s string) (VendorId, error) {
	s = strpredone(s)
	v, ok := list_vendor_name[s]
	if ok {
		return v, nil
	}
	return VENDOR_NO, ERR_VENDOR_INVALID
}

//对外提供重载vendor的方法 已存在的不允许重载
func Customer_SetVendorList(vendorid int, name string, typ string) error {
	id := VendorId(vendorid)
	if id == VENDOR_NO {
		return ERR_SET_VENDOR //不允许重载vendorid为0的
	}
	if _, ok := list_vendor_id[id]; ok { //若存在，则不允许重载
		return ERR_SET_VENDOR
	}
	name = strpredone(name)
	typ = strpredone(typ)
	if typ != "TYPE4" {
		typ = "IETF"
	}
	var c const_vendor = const_vendor{name, typ}
	list_vendor_id[id] = c
	list_vendor_name[name] = id
	return nil
}

//返回已加载的vendorid列表
func ListVenderId() []int {
	list := make([]int, 0)
	for id, _ := range list_vendor_id {
		list = append(list, int(id))
	}
	return list
}

//返回已加载的vendorname列表
func ListVenderName() []string {
	list := make([]string, 0)
	for ss, _ := range list_vendor_name {
		list = append(list, ss)
	}
	return list
}
