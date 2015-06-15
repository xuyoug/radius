package radius

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

//AttIdV定义厂商属性
type AttIdV struct {
	VendorId
	Id int
}

//ATTIDV_ERR定义错误的厂商属性
var ATTIDV_ERR AttIdV = AttIdV{VENDOR_ERR, 0}

//String方法返回AttIdV的字符串表达形式
func (a AttIdV) String() string {
	v, ok := list_attV_id[a.VendorId][a.Id]
	if ok {
		return a.VendorId.String() + ":" + v.Name + "(" + strconv.Itoa(a.Id) + ")"
	}
	return a.VendorId.String() + ":UNKNOWN_ATTRIBUTE(" + strconv.Itoa(a.Id) + ")"
}

//String方法返回AttIdV其标准名字
func (a AttIdV) AttributeName() string {
	v, ok := list_attV_id[a.VendorId][a.Id]
	if ok {
		return v.Name
	}
	return ""
}

//ValueType方法返回其值类型
func (a AttIdV) ValueType() string {
	v, ok := list_attV_id[a.VendorId][a.Id]
	if ok {
		return v.Type
	}
	return ""
}

//IsValid方法返回其是否有效
func (a AttIdV) IsValid() bool {
	_, ok := list_attV_id[a.VendorId][a.Id]
	return ok
}

//readAttIdV提供从buffer中读取AttIdV的方法
//发生错误则返回
func readAttIdV(buf *bytes.Buffer) (AttIdV, error) {
	var vid VendorId
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &vid)
	if !vid.IsValidVendor() { //不是有效vendor则返回错误
		return ATTIDV_ERR, ERR_VENDOR_INVALID
	}

	var vaid int
	if vid.IsType4() {
		var tmp uint32
		err := binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &tmp)
		if err != nil {
			return ATTIDV_ERR, err
		}
		vaid = int(tmp)
	} else {
		b, err := buf.ReadByte()
		if err != nil {
			return ATTIDV_ERR, err
		}
		vaid = int(b)
	}
	return AttIdV{vid, vaid}, nil //允许未知的属性
}

//Write方法将AttIdV自己写buffer
func (a AttIdV) write(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, a.VendorId)
	if a.IsType4() {
		binary.Write(buf, binary.BigEndian, uint32(a.Id))
	} else {
		binary.Write(buf, binary.BigEndian, uint8(a.Id))
	}
}

//getattV提供直接通过字符串获取厂商属性定义的方法
func getattidv(s string) AttIdV {
	for vid, v := range list_attV_name {
		for vaname, vaid := range v {
			if vaname == s {
				return AttIdV{vid, vaid}
			}
		}
	}
	return ATTIDV_ERR
}

//GetAttIdV提供根据字符串查找具体厂商属性的方法
//字符串以":"分隔
//":"之前为vendor名称，之后为属性名称
//若只有属性名称，则进行全部查找
func GetAttIdV(s string) AttIdV {
	s = stringfix(s)
	var vid VendorId
	ss := strings.Split(s, ":")
	if len(ss) == 1 {
		return getattidv(ss[0])
	}
	if len(ss) == 2 {
		vid = GetVendorId(ss[0])
		if vid != VENDOR_ERR {
			return ATTIDV_ERR
		}
		v, ok := list_attV_name[vid][ss[1]]
		if ok {
			return AttIdV{vid, v}
		}
	}
	return ATTIDV_ERR
}
