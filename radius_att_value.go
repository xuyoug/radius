package radius

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
)

//定义attributevalue的类型
type INTEGER uint32
type STRING string
type IPADDR []byte
type TAG_INT uint32
type TAG_STR string
type HEXADECIMAL []byte

//String格式化属性值输出
func (v INTEGER) String() string {
	return strconv.Itoa(int(v))
}
func (v STRING) String() string {
	return string(v)
}
func (v IPADDR) String() string {
	return fmt.Sprintf("%v", []byte(v))
}
func (v TAG_INT) String() string {
	return "TAG(" + strconv.Itoa(int(v>>24)) + ")" + strconv.Itoa(int(v))
}
func (v TAG_STR) String() string {
	tmp := []byte(v)
	return "TAG(" + string(tmp[0:1]) + ")" + string(v)
}
func (v HEXADECIMAL) String() string {
	return fmt.Sprintf("%v", []byte(v))
}

//ValueLen获取属性值的长度
func (v INTEGER) ValueLen() int {
	return 4
}
func (v STRING) ValueLen() int {
	return len(v)
}
func (v IPADDR) ValueLen() int {
	return len(v)
}
func (v TAG_INT) ValueLen() int {
	return 4
}
func (v TAG_STR) ValueLen() int {
	return len(v)
}
func (v HEXADECIMAL) ValueLen() int {
	return len(v)
}

//Value获取属性值的直接表达
func (v INTEGER) Value() interface{} {
	return int(v)
}
func (v STRING) Value() interface{} {
	return string(v)
}
func (v IPADDR) Value() interface{} {
	return []byte(v)
}
func (v TAG_INT) Value() interface{} {
	return int(v)
}
func (v TAG_STR) Value() interface{} {
	return string(v)
}
func (v HEXADECIMAL) Value() interface{} {
	return []byte(v)
}

//ValueTypestring获取属性值的类型字符串
func (v INTEGER) ValueTypestring() string {
	return "INTEGER"
}
func (v STRING) ValueTypestring() string {
	return "STRING"
}
func (v IPADDR) ValueTypestring() string {
	return "IPADDR"
}
func (v TAG_INT) ValueTypestring() string {
	return "TAG_INT"
}
func (v TAG_STR) ValueTypestring() string {
	return "TAG_STR"
}
func (v HEXADECIMAL) ValueTypestring() string {
	return "HEXADECIMAL"
}

//writeBuff写值入buffer
func (v INTEGER) writeBuffer(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v STRING) writeBuffer(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v IPADDR) writeBuffer(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v TAG_INT) writeBuffer(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v TAG_STR) writeBuffer(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v HEXADECIMAL) writeBuffer(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}

//AttributeValue定时属性值类型接口
type AttributeValue interface {
	writeBuffer(buf *bytes.Buffer)
	String() string
	Value() interface{}
	ValueLen() int
	ValueTypestring() string
}

//newAttributeValue获取空属性值
func NewAttributeValueEmpt(attrType string) (AttributeValue, error) {
	attrType = stringfix(attrType)
	switch attrType {
	case "INTEGER":
		var v INTEGER
		return v, nil
	case "STRING":
		var v STRING
		return v, nil
	case "IPADDR":
		var v IPADDR
		return v, nil
	case "TAG_INT":
		var v TAG_INT
		return v, nil
	case "TAG_STR":
		var v TAG_STR
		return v, nil
	case "HEXADECIMAL":
		var v HEXADECIMAL
		return v, nil
	default:
		return nil, ERR_ATTV_TYPE
	}
}

//newAttributeValueFromBuff从buff获取属性值
func newAttributeValueFromBuff(attrType string, length int, buf_in *bytes.Buffer) (AttributeValue, error) {
	v, err := NewAttributeValueEmpt(attrType)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(buf_in.Next(length))
	var tmp uint32
	switch attrType {
	case "INTEGER":
		if length != 4 {
			return nil, errors.New("INTEGER type must use length 4")
		}
		binary.Read(buf, binary.BigEndian, &tmp)
		v = INTEGER(tmp)
	case "STRING":
		v = STRING(buf.Bytes())
	case "IPADDR":
		s := string(buf.Bytes())
		p := net.ParseIP(s)
		if p != nil {
			v = IPADDR([]byte(p.To4()))
		} else {
			v = IPADDR(buf.Bytes())
		}
	case "TAG_INT":
		if length != 4 {
			return nil, errors.New("TAG_INT type must use length 4")
		}
		binary.Read(buf, binary.BigEndian, &tmp)
		v = TAG_INT(tmp)
	case "TAG_STR":
		v = TAG_STR(buf.Bytes())
	case "HEXADECIMAL":
		v = HEXADECIMAL(buf.Bytes())
	default:
		return nil, ERR_ATTV_TYPE
	}
	return v, nil
}

//NewAttributeValue根据指定内容生成属性值
//数字类型只接受int或uint32
//字符类型只接受string或[]byte
//该方法主要由外部使用
func NewAttributeValue(attrType string, i interface{}) (AttributeValue, error) {
	attrType = stringfix(attrType)
	v, err := NewAttributeValueEmpt(attrType)
	if err != nil {
		return INTEGER(0), err
	}
	switch attrType {
	case "INTEGER":
		tmp, ok := i.(int)
		if ok {
			v = INTEGER(tmp)
			return v, nil
		}
		tmp1, ok1 := i.(uint32)
		if ok1 {
			v = INTEGER(tmp1)
			return v, nil
		}
		return nil, errors.New("INTEGER type must fill by int or uint32 type")
	case "STRING":
		tmp, ok := i.(string)
		if ok {
			v = STRING(tmp)
			return v, nil
		}
		tmp1, ok1 := i.([]byte)
		if ok1 {
			v = STRING(tmp1)
			return v, nil
		}
		return nil, errors.New("STRING type must fill by string or []byte type")
	case "IPADDR":
		tmp, ok := i.(string)
		if ok {
			p := net.ParseIP(tmp)
			if p != nil {
				v = IPADDR([]byte(p.To4()))
			} else {
				v = IPADDR(tmp)
			}
			return v, nil
		}
		tmp1, ok1 := i.([]byte)
		if ok1 {
			v = IPADDR(tmp1)
			return v, nil
		}
		return nil, errors.New("IPADDR type must fill by string or []byte type")
	case "TAG_INT":
		tmp, ok := i.(int)
		if ok {
			v = TAG_INT(tmp)
			return v, nil
		}
		tmp1, ok1 := i.(uint32)
		if ok1 {
			v = TAG_INT(tmp1)
			return v, nil
		}
		return nil, errors.New("TAG_INT type must fill by int or uint32 type")
	case "TAG_STR":
		tmp, ok := i.(string)
		if ok {
			v = TAG_STR(tmp)
			return v, nil
		}
		tmp1, ok1 := i.([]byte)
		if ok1 {
			v = TAG_STR(tmp1)
			return v, nil
		}
		return nil, errors.New("TAG_STR type must fill by string or []byte type")
	case "HEXADECIMAL":
		tmp, ok := i.(string)
		if ok {
			v = HEXADECIMAL(tmp)
			return v, nil
		}
		tmp1, ok1 := i.([]byte)
		if ok1 {
			v = HEXADECIMAL(tmp1)
			return v, nil
		}
		return nil, errors.New("HEXADECIMAL type must fill by string or []byte type")
	default:
		return nil, ERR_ATTV_TYPE
	}
}

//NewAttributeValueS根据指定字符串生成属性值
//该方法主要由外部使用
func NewAttributeValueS(attrType, s string) (AttributeValue, error) {
	attrType = stringfix(attrType)
	v, err := NewAttributeValueEmpt(attrType)
	if err != nil {
		return INTEGER(0), err
	}
	switch attrType {
	case "INTEGER":
		i, err1 := strconv.Atoi(s)
		if err1 == nil {
			v = INTEGER(i)
			return v, nil
		}
		return nil, ERR_ATTV_TYPE
	case "STRING":
		v = STRING(s)
		return v, nil
	case "IPADDR":
		p := net.ParseIP(s)
		if p != nil {
			v = IPADDR([]byte(p.To4()))
		} else {
			v = IPADDR(s)
		}
		return v, nil
	case "TAG_INT":
		i, err1 := strconv.Atoi(s)
		if err1 == nil {
			v = INTEGER(i)
			return v, nil
		}
		return nil, ERR_ATTV_TYPE
	case "TAG_STR":
		v = TAG_STR(s)
		return v, nil
	case "HEXADECIMAL":
		v = HEXADECIMAL(s)
		return v, nil
	default:
		return nil, ERR_ATTV_TYPE
	}
}

//IsValidAttributeType判断是否是有效的属性值类型字符串
func IsValidAttributeValueType(s string) bool {
	if s == "INTEGER" || s == "STRING" || s == "IPADDR" || s == "TAG_INT" || s == "TAG_STR" || s == "HEXADECIMAL" {
		return true
	}
	return false
}
