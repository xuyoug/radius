package radius

import (
	//"fmt"
	"bytes"
	"encoding/binary"
	"strconv"
)

//定义厂商属性

//定义radius的attribute
type AttId uint8
type AttV uint8
type AttV4 uint32

func (a AttV) String() string {
	return strconv.Itoa(int(a))
}

func (a AttV4) String() string {
	return strconv.Itoa(int(a))
}

type AttVId struct {
	VendorId
	AttV
}

type AttV4Id struct {
	VendorId
	AttV4
}

func (a AttVId) String() string {
	return getAttVstring(a.VendorId, a.AttV)
}

func (a AttVId) Typestring() string {
	return getAttVtypestring(a.VendorId, a.AttV)
}

func (a AttVId) IsValid() bool {
	return isvalidAttV(a.VendorId, a.AttV)
}

func (a AttV4Id) String() string {
	return getAttV4string(a.VendorId, a.AttV4)
}

func (a AttV4Id) Typestring() string {
	return getAttV4typestring(a.VendorId, a.AttV4)
}

func (a AttV4Id) IsValid() bool {
	return isvalidAttV4(a.VendorId, a.AttV4)
}

//
func (a AttId) writevendor(buf *bytes.Buffer) {
}

func (a AttVId) writevendor(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, a.VendorId)
}

func (a AttV4Id) writevendor(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, a.VendorId)
}

//
func (a AttId) writeAtt(buf *bytes.Buffer) {
}

func (a AttVId) writeAtt(buf *bytes.Buffer) {
	buf.WriteByte(byte(a.AttV))
}

func (a AttV4Id) writeAtt(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, a.AttV4)
}

//定义属性标示符
type AttributeId interface {
	writevendor(buf *bytes.Buffer)
	writeAtt(buf *bytes.Buffer)
	String() string
	Typestring() string
	IsValid() bool
}

//NewAttributeId is a fucking and terrible thing!Fuck it!
//这种方式获取 时间在0.003ms-0.01ms之间
func NewAttributeId(in ...interface{}) (interface{}, error) {
	l := len(in)
	switch l {
	case 1:
		switch in[0].(type) {
		case int:
			if in[0].(int) != 26 {
				return AttId(in[0].(int)), nil
			}
			return nil, ERR_ATT_SET
		case string:
			var err error
			var aid AttId
			aid, err = GetAttId(in[0].(string))
			if err == nil {
				return aid, nil
			}
			var vid VendorId
			var va AttV
			vid, va, err = GetAttV(in[0].(string))
			if err == nil {
				return AttVId{vid, va}, nil
			}
			var va4 AttV4
			vid, va4, err = GetAttV4(in[0].(string))
			if err == nil {
				return AttV4Id{vid, va4}, nil
			}
			return nil, ERR_ATT_SET

		case AttId:
			return in[0].(AttId), nil

		default:
			return nil, ERR_ATT_SET
		}
	case 2:
		switch in[0].(type) {
		case int:
			i := in[0].(int)
			switch in[1].(type) {
			case int:
				j := in[1].(int)
				if !VendorId(i).IsvalidVendor() {
					return nil, ERR_ATT_SET
				}
				if i == 0 {
					return AttId(j), nil
				}
				if VendorId(i).Typestring() == "IETF" {
					return AttVId{VendorId(i), AttV(j)}, nil
				}
				if VendorId(i).Typestring() == "TYPE4" {
					return AttV4Id{VendorId(i), AttV4(j)}, nil
				}
				return nil, ERR_ATT_SET
			case AttV:
				if !VendorId(i).IsvalidVendor() || i == 0 {
					return nil, ERR_ATT_SET
				}
				if VendorId(i).Typestring() == "IETF" {
					return AttVId{VendorId(i), in[1].(AttV)}, nil
				}
				if VendorId(i).Typestring() == "TYPE4" {
					return nil, ERR_ATT_SET
				}
				return nil, ERR_ATT_SET
			case AttV4:
				if !VendorId(i).IsvalidVendor() || i == 0 {
					return nil, ERR_ATT_SET
				}
				if VendorId(i).Typestring() == "IETF" {
					return nil, ERR_ATT_SET
				}
				if VendorId(i).Typestring() == "TYPE4" {
					return AttV4Id{VendorId(i), in[1].(AttV4)}, nil
				}
				return nil, ERR_ATT_SET
			case string:
				if !VendorId(i).IsvalidVendor() {
					return nil, ERR_ATT_SET
				}
				if i == 0 {
					vid, err := GetAttId(in[1].(string))
					if err == nil {
						return vid, nil
					}
					return nil, ERR_ATT_SET
				}
				if VendorId(i).Typestring() == "IETF" {
					vaid, ok := list_attv_name[VendorId(i)][in[1].(string)]
					if ok {
						return AttVId{VendorId(i), vaid}, nil
					}
				}
				if VendorId(i).Typestring() == "TYPE4" {
					vaid4, ok := list_attv4_name[VendorId(i)][in[1].(string)]
					if ok {
						return AttV4Id{VendorId(i), vaid4}, nil
					}
				}
				return nil, ERR_ATT_SET
			default:
				return nil, ERR_ATT_SET
			}
		case VendorId:
			vid := in[0].(VendorId)
			if !vid.IsvalidVendor() {
				return nil, ERR_ATT_SET
			}
			switch in[1].(type) {
			case int:
				if vid == VENDOR_NO {
					return AttId(in[1].(int)), nil
				}
				if vid.Typestring() == "IETF" {
					return AttVId{vid, AttV(in[1].(int))}, nil
				}
				if vid.Typestring() == "TYPE4" {
					return AttV4Id{vid, AttV4(in[1].(int))}, nil
				}
				return nil, ERR_ATT_SET
			case AttV:
				if vid == VENDOR_NO {
					return nil, ERR_ATT_SET
				}
				if vid.Typestring() == "IETF" {
					return AttVId{vid, in[1].(AttV)}, nil
				}
				if vid.Typestring() == "TYPE4" {
					return nil, ERR_ATT_SET
				}
				return nil, ERR_ATT_SET
			case AttV4:
				if vid == VENDOR_NO {
					return nil, ERR_ATT_SET
				}
				if vid.Typestring() == "IETF" {
					return nil, ERR_ATT_SET
				}
				if vid.Typestring() == "TYPE4" {
					return AttV4Id{vid, in[1].(AttV4)}, nil
				}
				return nil, ERR_ATT_SET
			case string:
				if vid == VENDOR_NO {
					return GetAttId(in[1].(string))
				}
				if vid.Typestring() == "IETF" {
					vaid, ok := list_attv_name[vid][in[1].(string)]
					if ok {
						return AttVId{vid, vaid}, nil
					}
				}
				if vid.Typestring() == "TYPE4" {
					vaid4, ok := list_attv4_name[vid][in[1].(string)]
					if ok {
						return AttV4Id{vid, vaid4}, nil
					}
				}
				return nil, ERR_ATT_SET
			default:
				return nil, ERR_ATT_SET
			}
		case string:
			vid, err := GetVendorId(in[0].(string))
			if err != nil {
				return nil, err
			}
			if !vid.IsvalidVendor() {
				return nil, ERR_ATT_SET
			}
			switch in[1].(type) {
			case int:
				if vid == VENDOR_NO {
					return AttId(in[1].(int)), nil
				}
				if vid.Typestring() == "IETF" {
					return AttVId{vid, AttV(in[1].(int))}, nil
				}
				if vid.Typestring() == "TYPE4" {
					return AttV4Id{vid, AttV4(in[1].(int))}, nil
				}
				return nil, ERR_ATT_SET
			case AttV:
				if vid == VENDOR_NO {
					return nil, ERR_ATT_SET
				}
				if vid.Typestring() == "IETF" {
					return AttVId{vid, in[1].(AttV)}, nil
				}
				if vid.Typestring() == "TYPE4" {
					return nil, ERR_ATT_SET
				}
				return nil, ERR_ATT_SET
			case AttV4:
				if vid == VENDOR_NO {
					return nil, ERR_ATT_SET
				}
				if vid.Typestring() == "IETF" {
					return nil, ERR_ATT_SET
				}
				if vid.Typestring() == "TYPE4" {
					return AttV4Id{vid, in[1].(AttV4)}, nil
				}
				return nil, ERR_ATT_SET
			case string:
				if vid == VENDOR_NO {
					return GetAttId(in[1].(string))
				}
				if vid.Typestring() == "IETF" {
					vaid, ok := list_attv_name[vid][in[1].(string)]
					if ok {
						return AttVId{vid, vaid}, nil
					}
				}
				if vid.Typestring() == "TYPE4" {
					vaid4, ok := list_attv4_name[vid][in[1].(string)]
					if ok {
						return AttV4Id{vid, vaid4}, nil
					}
				}
				return nil, ERR_ATT_SET
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

//定义属性的长度
type attributeLen uint8

//

// AttrList
type AttributeList struct {
	list_name  map[AttributeId][]AttributeValue
	list_order []AttributeId
}

func (a *AttributeList) AddAttr(r AttributeId, v AttributeValue) {
	if _, ok := a.list_name[r]; !ok {
		a.list_name[r] = make([]AttributeValue, 0)
		a.list_name[r] = append(a.list_name[r], v)
		a.list_order = append(a.list_order, r)
	} else {
		a.list_name[r] = append(a.list_name[r], v)
	}
}

func (a *AttributeList) Len() int {
	return len(a.list_order)
}

func (a *AttributeList) GetAttrs() []AttributeId {
	return a.list_order
}

func (a *AttributeList) GetAttr(r AttributeId) ([]AttributeValue, error) {
	if v, ok := a.list_name[r]; ok {
		return v, nil
	}
	return nil, ERR_ATT_NO
}

func (a *AttributeList) String() string {
	var s string
	s += "Attributes:\n"
	for i, v := range a.list_order {
		for _, vv := range a.list_name[v] {
			s += strconv.Itoa(i+1) + ": " + v.String() + " value:" + vv.String() + "\n"
		}
	}
	return s
}
