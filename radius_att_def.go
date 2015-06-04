package radius

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

//AttId定义标准属性Id
type AttId uint8

//String方法返回AttId其名字
func (a AttId) String() string {
	s, ok := list_attributestand_id[a]
	if ok {
		return s.Name
	}
	return ""
}

//String方法返回AttId其标准名字
func (a AttId) AttributeName() string {
	s, ok := list_attributestand_id[a]
	if ok {
		return s.Name
	}
	return ""
}

//ValueTypestring方法返回其值类型字符串类型
func (a AttId) ValueTypestring() string {
	s, ok := list_attributestand_id[a]
	if ok {
		return s.Type
	}
	return ""
}

//Typestring方法返回AttId是否有效
func (a AttId) IsValid() bool {
	_, ok := list_attributestand_id[a]
	return ok
}

//Write方法将AttId自己写入buffer
func (a AttId) writeAttributeId(buf *bytes.Buffer) error {
	err := buf.WriteByte(byte(a))
	if err != nil {
		return err
	}
	return nil
}

//readAttId提供从buffer读取AttId的方法
func readAttId(buf *bytes.Buffer) (AttId, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return ATT_NO, err
	}
	return AttId(b), nil
}

//GetAttId提供根据名字返回AttId的方法
func GetAttId(s string) (AttId, error) {
	s = stringfix(s)
	a, ok := list_attributestand_name[s]
	if ok {
		return a, nil
	}
	return ATT_NO, ERR_ATT_UNK
}

//AttIdV定义厂商属性
type AttIdV struct {
	VendorId
	Id int
}

//ATTIDV_ERR定义错误的厂商属性
var ATTIDV_ERR AttIdV = AttIdV{VENDOR_NO, 0}

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

//ValueTypestringTypestring方法返回其值类型
func (a AttIdV) ValueTypestring() string {
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
	vtype := vid.VendorTypestring()
	var vaid int
	if vtype == "IETF" {
		b, err := buf.ReadByte()
		if err != nil {
			return ATTIDV_ERR, ERR_RADIUS_FMT
		}
		vaid = int(b)
	}
	if vtype == "TYPE4" {
		var tmp uint32
		binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &tmp)
		vaid = int(tmp)
	}
	return AttIdV{vid, vaid}, nil
}

//WriteAttributeId方法将AttIdV自己写buffer
func (a AttIdV) writeAttributeId(buf *bytes.Buffer) error {
	err := binary.Write(buf, binary.BigEndian, a.VendorId)
	if err != nil {
		return err
	}
	typ := a.VendorTypestring()
	if typ == "IETF" {
		err = binary.Write(buf, binary.BigEndian, uint8(a.Id))
		if err != nil {
			return err
		}
	}
	if typ == "TYPE4" {
		err = binary.Write(buf, binary.BigEndian, uint32(a.Id))
		if err != nil {
			return err
		}
	}
	return nil
}

//getattV提供直接通过字符串获取厂商属性定义的方法
func getattidv(s string) (AttIdV, error) {
	for vid, v := range list_attV_name {
		for vaname, vaid := range v {
			if vaname == s {
				return AttIdV{vid, vaid}, nil
			}
		}
	}
	return ATTIDV_ERR, ERR_ATT_UNK
}

//GetAttIdV提供根据字符串查找具体厂商属性的方法
//字符串以":"分隔
//":"之前为vendor名称，之后为属性名称
//若只有属性名称，则进行全部查找
func GetAttIdV(s string) (AttIdV, error) {
	s = stringfix(s)
	var vid VendorId
	var err error
	ss := strings.Split(s, ":")
	if len(ss) == 1 {
		return getattidv(ss[0])
	}
	if len(ss) == 2 {
		vid, err = GetVendorId(ss[0])
		if err != nil {
			return ATTIDV_ERR, err
		}
		v, ok := list_attV_name[vid][ss[1]]
		if ok {
			return AttIdV{vid, v}, nil
		}
	}
	return ATTIDV_ERR, ERR_ATT_UNK
}

//定义属性标示符接口
type AttributeId interface {
	writeAttributeId(buf *bytes.Buffer) error
	String() string
	ValueTypestring() string
	IsValid() bool
	AttributeName() string
}

//NewAttributeId is a fucking and terrible thing!Fuck it!
//这种方式获取 时间在0.003ms-0.01ms之间
func NewAttributeId(in ...interface{}) (AttributeId, error) {
	var vid VendorId
	var aid AttId
	var aidv AttIdV
	var err error

	l := len(in)
	switch l {
	case 1:
		switch in[0].(type) {
		case int:
			if in[0].(int) != 26 || in[0].(int) < 255 {
				return AttId(in[0].(int)), nil
			}
			return ATT_NO, ERR_ATT_UNK
		case string:
			aid, err = GetAttId(in[0].(string))
			if err == nil {
				return aid, nil
			}
			aidv, err = GetAttIdV(in[0].(string))
			if err == nil {
				return aidv, nil
			}
			return ATT_NO, ERR_ATT_UNK
		case AttId:
			aid = in[0].(AttId)
			if aid == ATTID_VENDOR_SPECIFIC {
				return ATT_NO, ERR_ATT_SET
			}
			return aid, nil
		default:
			return ATT_NO, ERR_ATT_SET
		}
	case 2:
		switch in[0].(type) {
		case int:
			vid = VendorId(in[0].(int))
			if !vid.IsValidVendor() {
				return ATT_NO, ERR_ATT_SET
			}
			switch in[1].(type) {
			case int:
				aidv, err = vid.GetAttById(in[1].(int))
				if err == nil {
					return aidv, nil
				}
				return ATT_NO, err
			case string:
				aidv, err = vid.GetAttByName(in[1].(string))
				if err == nil {
					return aidv, nil
				}
				return ATT_NO, err
			default:
				return ATT_NO, ERR_ATT_SET
			}
		case VendorId:
			vid = in[0].(VendorId)
			if !vid.IsValidVendor() {
				return nil, ERR_ATT_SET
			}
			switch in[1].(type) {
			case int:
				aidv, err = vid.GetAttById(in[1].(int))
				if err == nil {
					return aidv, nil
				}
				return ATT_NO, err
			case string:
				aidv, err = vid.GetAttByName(in[1].(string))
				if err == nil {
					return aidv, nil
				}
				return ATT_NO, err
			default:
				return ATT_NO, ERR_ATT_SET
			}
		case string:
			vid, err = GetVendorId(in[0].(string))
			if err != nil {
				return nil, err
			}
			switch in[1].(type) {
			case int:
				aidv, err = vid.GetAttById(in[1].(int))
				if err == nil {
					return aidv, nil
				}
				return ATT_NO, err
			case string:
				aidv, err = vid.GetAttByName(in[1].(string))
				if err == nil {
					return aidv, nil
				}
				return ATT_NO, err
			default:
				return ATT_NO, ERR_ATT_SET
			}
		default:
			return nil, ERR_ATT_SET
		}
	default:
		return nil, ERR_ATT_SET
	}
	return nil, ERR_ATT_SET
}

//Attribute定义一个完整的属性，包含属性描述及属性值
type Attribute struct {
	AttributeId
	AttributeValue
}

//AttributeMeanig返回属性表示的含义
func (a *Attribute) AttributeMeanig() string {
	s := a.AttributeName()
	if a.AttributeValue.ValueTypestring() != "INTEGER" && a.AttributeValue.ValueTypestring() != "TAG_INT" {
		return ""
	}
	v := a.Value().(int)
	m, ok := list_attribute_meaning[s][uint32(v)]
	if ok {
		return m
	}
	return ""
}

//AttributeMeanig返回属性表示的含义
func (a *Attribute) String() string {
	var s string
	s += a.AttributeId.String() + " value:" + a.AttributeValue.String()
	if meaning := a.AttributeMeanig(); meaning != "" {
		s += " ("
		s += meaning
		s += ")"
	}
	return s
}

//ATTIDV_ERR定义错误的厂商属性
var ATTRIBUTE_ERR Attribute = Attribute{ATT_NO, INTEGER(0)}

//
func readAttribute(buf *bytes.Buffer) (Attribute, error) {
	var attid AttId
	var length, lengthv int
	var b byte
	var err error
	var typ, vtyp string
	attid, err = readAttId(buf)
	if err != nil {
		return ATTRIBUTE_ERR, err
	}
	b, err = buf.ReadByte()
	if err != nil {
		return ATTRIBUTE_ERR, ERR_ATT_FMT
	}
	length = int(b)
	if attid != ATTID_VENDOR_SPECIFIC {
		typ = attid.ValueTypestring()
		v, err1 := newAttributeValueFromBuff(typ, length-2, buf)
		if err1 != nil {
			return ATTRIBUTE_ERR, err1
		}
		return Attribute{attid, v}, nil
	} else {
		attidv, err1 := readAttIdV(buf)
		if err1 != nil {
			return ATTRIBUTE_ERR, err1
		}
		vtyp = attidv.VendorTypestring()
		typ = attidv.ValueTypestring()
		if vtyp == "IETF" {
			b, err = buf.ReadByte()
			if err != nil {
				return ATTRIBUTE_ERR, ERR_ATT_FMT
			}
			lengthv = int(b)
			if lengthv != length-6 {
				return ATTRIBUTE_ERR, ERR_ATT_FMT
			}
			v, err1 := newAttributeValueFromBuff(typ, lengthv-2, buf)
			if err1 != nil {
				return ATTRIBUTE_ERR, err1
			}
			return Attribute{attidv, v}, nil
		}
		if vtyp == "TYPE4" {
			v, err1 := newAttributeValueFromBuff(typ, length-10, buf)
			if err1 != nil {
				return ATTRIBUTE_ERR, err1
			}
			return Attribute{attidv, v}, nil
		}
	}
	return ATTRIBUTE_ERR, ERR_ATT_FMT
}

//writebuf方法将属性写buffer
func (v *Attribute) writeBuffer(buf *bytes.Buffer) {
	switch v.AttributeId.(type) {
	case AttId:
		v.AttributeId.writeAttributeId(buf)
		buf.WriteByte(byte(uint8(v.AttributeValue.ValueLen() + 2)))
		v.AttributeValue.writeBuffer(buf)
	case AttIdV:
		buf.WriteByte(byte(ATTID_VENDOR_SPECIFIC))
		if v.AttributeId.(AttIdV).VendorTypestring() == "IETF" {
			buf.WriteByte(byte(uint8(v.AttributeValue.ValueLen() + 8)))
			v.AttributeId.writeAttributeId(buf)
			buf.WriteByte(byte(uint8(v.AttributeValue.ValueLen() + 2)))
		}
		if v.AttributeId.(AttIdV).VendorTypestring() == "TYPE4" {
			buf.WriteByte(byte(uint8(v.AttributeValue.ValueLen() + 10)))
			v.AttributeId.writeAttributeId(buf)
		}
		v.AttributeValue.writeBuffer(buf)
	}
}

// AttrList
type AttributeList struct {
	attributes []Attribute
}

//
func (a *AttributeList) AddAttr(r Attribute) {
	a.attributes = append(a.attributes, r)
}

//
func (a *AttributeList) GetAttrsNum() int {
	return len(a.attributes)
}

//
func (a *AttributeList) GetAttrs() ([]Attribute, int) {
	return a.attributes, len(a.attributes)
}

//
func (a *AttributeList) GetAttr(r AttributeId) ([]AttributeValue, int) {
	list := make([]AttributeValue, 0)
	var numbers int
	for _, v := range a.attributes {
		if v.AttributeId == r {
			list = append(list, v.AttributeValue)
			numbers += 1
		}
	}
	return list, numbers
}

//
func (a *AttributeList) GetAttrFist(r AttributeId) (AttributeValue, error) {
	for _, v := range a.attributes {
		if v.AttributeId == r {
			return v.AttributeValue, nil
		}
	}
	return INTEGER(0), ERR_ATT_NO
}

//
func (a *AttributeList) String() string {
	var s string
	s += "Attributes:"
	for i, v := range a.attributes {
		s += "\n"
		s += strconv.Itoa(i+1) + ": " + v.String()
	}
	return s
}
