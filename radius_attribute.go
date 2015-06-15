package radius

import (
	"bytes"
)

//定义属性标示符接口
type AttributeId interface {
	write(buf *bytes.Buffer) error
	String() string
	ValueType() string
	IsValid() bool
	AttributeName() string
}

//newAttributeId is a fucking and terrible thing!Fuck it!
//这种方式获取 时间在0.003ms-0.01ms之间
func newAttributeId(in ...interface{}) (AttributeId, error) {
	var vid VendorId
	var aid AttId
	var aidv AttIdV
	var err error

	switch len(in) {
	case 1:
		switch in[0].(type) {
		case int:
			if _, ok := list_attributestand_id[AttId(in[0].(int))]; ok && in[0].(int) != 26 {
				return AttId(in[0].(int)), nil
			}
			return nil, ERR_ATT_UKN
		case string:
			aid = GetAttId(in[0].(string))
			if aid != ATTID_ERR {
				return aid, nil
			}
			aidv = GetAttIdV(in[0].(string))
			if aidv != ATTIDV_ERR {
				return aidv, nil
			}
			return nil, ERR_ATT_UKN
		case AttId: //注意不要直接转了传进来，这里不检查错误
			aid = in[0].(AttId)
			if aid == ATTID_VENDOR_SPECIFIC {
				return nil, ERR_ATT_SET
			}
			return aid, nil
		default:
			return nil, ERR_ATT_SET
		}
	case 2:
		switch in[0].(type) {
		case int:
			vid = VendorId(in[0].(int))
			if !vid.IsValidVendor() {
				return nil, ERR_ATT_SET
			}
			switch in[1].(type) {
			case int:
				return vid.GetAttById(in[1].(int))
			case string:
				return vid.GetAttByName(in[1].(string))
			default:
				return nil, ERR_ATT_SET
			}
		case VendorId:
			vid = in[0].(VendorId)
			if !vid.IsValidVendor() {
				return nil, ERR_ATT_SET
			}
			switch in[1].(type) {
			case int:
				return vid.GetAttById(in[1].(int))
			case string:
				return vid.GetAttByName(in[1].(string))
			default:
				return nil, ERR_ATT_SET
			}
		case string:
			vid, err = GetVendorId(in[0].(string))
			if err != nil {
				return nil, err
			}
			switch in[1].(type) {
			case int:
				return vid.GetAttById(in[1].(int))
			case string:
				return vid.GetAttByName(in[1].(string))
			default:
				return nil, ERR_ATT_SET
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
	if a.ValueLen() != 4 {
		return ""
	}
	var i uint32
	switch a.AttributeValue.(type) {
	case INTEGER:
		i = uint32(a.AttributeValue.(INTEGER))
	case TAG_INT:
		i = uint32(a.AttributeValue.(TAG_INT))
	case IPADDR:
		i = a.AttributeValue.(IPADDR).Uint32()
	}
	return AttributeMeanig(s, i)
}

//重构String方法
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

//readAttribute从buf读取Attribute
func readAttribute(buf *bytes.Buffer) (Attribute, error) {
	var attid AttId
	var length, lengthv int
	var b byte
	var err error
	var typ string
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
		v, err1 := readAttributeValue(attid.ValueType(), length-2, buf)
		if err1 != nil {
			return ATTRIBUTE_ERR, err1
		}
		return Attribute{attid, v}, nil
	} else {
		attidv, err1 := readAttIdV(buf)
		if err1 != nil {
			return ATTRIBUTE_ERR, err1
		}
		typ = attidv.ValueTypestring()
		if !attidv.IsType4() {
			b, err = buf.ReadByte()
			lengthv = int(b)
			if lengthv != length-6 {
				return ATTRIBUTE_ERR, ERR_ATT_FMT
			}
			v, err1 := readAttributeValue(typ, lengthv-2, buf)
			if err1 != nil {
				return ATTRIBUTE_ERR, err1
			}
			return Attribute{attidv, v}, nil
		} else {
			v, err1 := readAttributeValue(typ, length-10, buf)
			if err1 != nil {
				return ATTRIBUTE_ERR, err1
			}
			return Attribute{attidv, v}, nil
		}
	}
	return ATTRIBUTE_ERR, ERR_ATT_FMT
}

//writebuf方法将属性写buffer
func (v *Attribute) write(buf *bytes.Buffer) {
	switch v.AttributeId.(type) {
	case AttId:
		v.AttributeId.writeAttributeId(buf)
		buf.WriteByte(byte(uint8(v.ValueLen() + 2)))
		v.AttributeValue.writeBuffer(buf)
	case AttIdV:
		buf.WriteByte(byte(ATTID_VENDOR_SPECIFIC))
		if !v.AttributeId.(AttIdV).IsType4() {
			buf.WriteByte(byte(uint8(v.ValueLen() + 8)))
			v.AttributeId.writeAttributeId(buf)
			buf.WriteByte(byte(uint8(v.ValueLen() + 2)))
		} else {
			buf.WriteByte(byte(uint8(v.ValueLen() + 10)))
			v.AttributeId.writeAttributeId(buf)
		}
		v.AttributeValue.writeBuffer(buf)
	}
}
