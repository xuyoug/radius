package radius

//定义常用厂商列表
const (
	VENDOR_ERR        VendorId = 0 //本包中  此vendorid 表示错误
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

//定义VEBDOR的ID映射
var list_vendor_id map[VendorId]Description_definition = map[VendorId]Description_definition{
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

//定义TYPE4的vendorlist
var list_vendor_type4 []VendorId = []VendorId{VENDOR_3COM}

//Customer_SetVendorList对外提供重载vendor的方法
//已存在的不允许重载
//vendortype不为TYPE4的均设置为IETF
func Customer_SetVendorList(vendorid int, name string, typ string) error {
	id := VendorId(vendorid)
	if id == VENDOR_ERR {
		return ERR_VENDOR_SET //不允许设置vendorid为0的
	}
	if _, ok := list_vendor_id[id]; ok { //若存在，则不允许重载
		return ERR_VENDOR_SET
	}
	name = stringfix(name)
	typ = stringfix(typ)
	if typ != "TYPE4" {
		typ = "IETF"
	}
	var c Description_definition = Description_definition{name, typ}
	list_vendor_id[id] = c
	list_vendor_name[name] = id
	if typ == "TYPE4" {
		list_vendor_type4 = append(list_vendor_type4, id)
	}
	return nil
}

//ListVender返回已加载的vendorid列表
func ListVender() map[VendorId]Description_definition {
	return list_vendor_id
}
