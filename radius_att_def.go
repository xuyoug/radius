package radius

import (
	//"fmt"
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

//定义标准属性列表
type AttId uint8

//String返回AttId其名字
func (a AttId) String() string {
	s, ok := list_attributestand_id[a]
	if ok {
		return s.Name
	}
	return ""
}

//Typestring返回AttId其类型
func (a AttId) Typestring() string {
	s, ok := list_attributestand_id[a]
	if ok {
		return s.Type
	}
	return ""
}

//Typestring返回AttId其类型
func (a AttId) IsValid() bool {
	_, ok := list_attributestand_id[a]
	return ok
}

func (a AttId) Write(buf *bytes.Buffer) error {
	err := buf.WriteByte(byte(a))
	if err != nil {
		return err
	}
	return nil
}

func readAttId(buf *bytes.Buffer) (AttId, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return ATT_NO, err
	}
	return AttId(b), nil
}

//根据名字返回AttId
func GetAttId(s string) (AttId, error) {
	s = stringfix(s)
	a, ok := list_attributestand_name[s]
	if ok {
		return a, nil
	}
	return ATT_NO, ERR_ATT_UNK
}

//定义厂商属性

//定义radius的attribute
type AttV uint8
type AttV4 uint32

func (a AttV) bytes() []byte {
	bs := make([]byte, 1)
	bs[0] = byte(a)
	return bs
}

func (a AttV4) bytes() []byte {
	// bs := make([]byte,4)
	// for i:=0;i<4;i++{
	// 	tmp := a
	// 	tmp<<i*8
	// 	tmp>>(3-i)*8
	// 	b := byte(tmp)
	// 	bs[i]=b
	// }
	// return bs
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, a)
	return buf.Bytes()
}

type AttVS interface {
	bytes() []byte
}

//厂商属性定义
type AttIdV struct {
	VendorId
	AttVS
}

var ATTIDV_ERR AttIdV = AttIdV{VENDOR_NO, nil}

func (a AttIdV) String() string {
	if a.VendorId.Typestring() == "IETF" {
		v, ok := list_AttV_id[a.VendorId][a.AttVS.(AttV)]
		if ok {
			return a.VendorId.String() + ":" + v.Name + "(" + strconv.Itoa(int(a.AttVS.(AttV))) + ")"
		}
		return a.VendorId.String() + ":UNKNOWN_ATTRIBUTE(" + strconv.Itoa(int(a.AttVS.(AttV))) + ")"
	}
	if a.VendorId.Typestring() == "TYPE4" {
		v, ok := list_AttV4_id[a.VendorId][a.AttVS.(AttV4)]
		if ok {
			return a.VendorId.String() + ":" + v.Name + "(" + strconv.Itoa(int(a.AttVS.(AttV4))) + ")"
		}
		return a.VendorId.String() + ":UNKNOWN_ATTRIBUTE(" + strconv.Itoa(int(a.AttVS.(AttV4))) + ")"
	}
	return ""
}

func (a AttIdV) Typestring() string {
	if a.VendorId.Typestring() == "IETF" {
		v, ok := list_AttV_id[a.VendorId][a.AttVS.(AttV)]
		if ok {
			return v.Type
		}
	}
	if a.VendorId.Typestring() == "TYPE4" {
		v, ok := list_AttV4_id[a.VendorId][a.AttVS.(AttV4)]
		if ok {
			return v.Type
		}
	}
	return ""
}

func (a AttIdV) IsValid() bool {
	if a.VendorId.Typestring() == "IETF" {
		_, ok := list_AttV_id[a.VendorId][a.AttVS.(AttV)]
		return ok
	}
	if a.VendorId.Typestring() == "TYPE4" {
		_, ok := list_AttV4_id[a.VendorId][a.AttVS.(AttV4)]
		return ok
	}
	return false
}

//getAttV直接通过字符串获取  不推荐
func getAttV(s string) (AttIdV, error) {
	for vid, v := range list_AttV_name {
		for vaname, vaid := range v {
			if vaname == s {
				return AttIdV{vid, vaid}, nil
			}
		}
	}
	for vid, v := range list_AttV4_name {
		for vaname, vaid := range v {
			if vaname == s {
				return AttIdV{vid, vaid}, nil
			}
		}
	}
	return ATTIDV_ERR, ERR_ATT_UNK
}

//
func GetAttV(s string) (AttIdV, error) {
	s = stringfix(s)
	var vid VendorId
	var err error
	ss := strings.Split(s, ":")
	if len(ss) == 1 {
		return getAttV(ss[0])
	}
	if len(ss) == 2 {
		vid, err = GetVendorId(ss[0])
		if err != nil {
			return ATTIDV_ERR, err
		}
		//根据vendorid判断vaid的类型
		switch vid.Typestring() {
		case "IETF":
			v, ok := list_AttV_name[vid][ss[1]]
			if ok {
				return AttIdV{vid, v}, nil
			}
		case "TYPE4":
			v, ok := list_AttV4_name[vid][ss[1]]
			if ok {
				return AttIdV{vid, v}, nil
			}
		}
	}
	return ATTIDV_ERR, ERR_ATT_UNK
}

//
func readAttIdV(buf *bytes.Buffer) (AttIdV, error) {
	var vid VendorId
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &vid)
	if !vid.IsValidVendor() {
		return ATTIDV_ERR, ERR_VENDOR_INVALID
	}
	var attidvs AttIdV
	vtype := vid.Typestring()
	if vtype == "IETF" {
		var vaid AttV
		binary.Read(bytes.NewBuffer(buf.Next(1)), binary.BigEndian, &vaid)
		attidvs = AttIdV{vid, vaid}
	}
	if vtype == "TYPE4" {
		var vaid AttV4
		binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &vaid)
		attidvs = AttIdV{vid, vaid}
	}
	return attidvs, nil
}

func (a AttIdV) Write(buf *bytes.Buffer) error {
	binary.Write(buf, binary.BigEndian, a.VendorId)
	_, err := buf.Write(a.AttVS.bytes())
	if err != nil {
		return err
	}
	return nil
}

//定义属性标示符
type AttributeId interface {
	Write(buf *bytes.Buffer) error
	String() string
	Typestring() string
	IsValid() bool
}

//NewAttributeId is a fucking and terrible thing!Fuck it!
//这种方式获取 时间在0.003ms-0.01ms之间
func NewAttributeId(in ...interface{}) (interface{}, error) {
	var vid VendorId
	var aid AttId
	var aidv AttIdV
	var err error

	l := len(in)
	switch l {
	case 1:
		switch in[0].(type) {
		case int:
			if in[0].(int) != 26 || in[0].(int) >= 255 {
				return AttId(in[0].(int)), nil
			}
			return nil, ERR_ATT_SET
		case string:
			aid, err = GetAttId(in[0].(string))
			if err == nil {
				return aid, nil
			}
			aidv, err = GetAttV(in[0].(string))
			if err == nil {
				return aidv, nil
			}
			return nil, ERR_ATT_SET
		case AttId:
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
				aidv, err = vid.getAttById(in[1].(int))
				if err == nil {
					return aidv, nil
				}
				return nil, err
			case string:
				aidv, err = vid.getAttByName(in[1].(string))
				if err == nil {
					return aidv, nil
				}
				return nil, err
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
				aidv, err = vid.getAttById(in[1].(int))
				if err == nil {
					return aidv, nil
				}
				return nil, err
			case string:
				aidv, err = vid.getAttByName(in[1].(string))
				if err == nil {
					return aidv, nil
				}
				return nil, err
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
				aidv, err = vid.getAttById(in[1].(int))
				if err == nil {
					return aidv, nil
				}
				return nil, err
			case string:
				aidv, err = vid.getAttByName(in[1].(string))
				if err == nil {
					return aidv, nil
				}
				return nil, err
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

func (a *AttributeList) GetAttrsNum() int {
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

func (a *AttributeList) GetAttrFist(r AttributeId) (AttributeValue, error) {
	if v, ok := a.list_name[r]; ok {
		return v[0], nil
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
