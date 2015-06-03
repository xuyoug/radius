package radius

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

//定义attributevalue的类型
type INTEGER uint32
type STRING string
type IPADDR []byte
type TAG_INT uint32
type TAG_STR string
type HEXADECIMAL []byte

//格式化输出
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

//获取长度
func (v INTEGER) Len() uint8 {
	return 4
}
func (v STRING) Len() uint8 {
	return uint8(len(v))
}
func (v IPADDR) Len() uint8 {
	return uint8(len(v))
}
func (v TAG_INT) Len() uint8 {
	return 4
}
func (v TAG_STR) Len() uint8 {
	return uint8(len(v))
}
func (v HEXADECIMAL) Len() uint8 {
	return uint8(len(v))
}

//获取值的直接表达
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

//获取类型字符串
func (v INTEGER) Typestring() string {
	return "INTEGER"
}
func (v STRING) Typestring() string {
	return "STRING"
}
func (v IPADDR) Typestring() string {
	return "IPADDR"
}
func (v TAG_INT) Typestring() string {
	return "TAG_INT"
}
func (v TAG_STR) Typestring() string {
	return "TAG_STR"
}
func (v HEXADECIMAL) Typestring() string {
	return "HEXADECIMAL"
}

//写值入buffer
func (v INTEGER) writeBuff(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v STRING) writeBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v IPADDR) writeBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v TAG_INT) writeBuff(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v TAG_STR) writeBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v HEXADECIMAL) writeBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}

//定时属性值类型接口
type AttributeValue interface {
	writeBuff(buf *bytes.Buffer)
	String() string
	Value() interface{}
	Len() uint8
	Typestring() string
}

//获取空属性值
func newAttributeValue(attrType string) (AttributeValue, error) {
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
		return nil, ERR_ATT_TYPE
	}
}

//从buff获取属性值
func newAttributeValueFromBuff(attrType string, length int, buf_in *bytes.Buffer) (AttributeValue, error) {
	v, err := newAttributeValue(attrType)
	if err != nil {
		return INTEGER(0), err
	}
	buf := bytes.NewBuffer(buf_in.Next(length))
	switch attrType {
	case "INTEGER":
		if length != 4 {
			return nil, ERR_ATT_TYPE
		}
		binary.Read(buf, binary.BigEndian, &v)
	case "STRING":
		v = STRING(buf.Bytes())
	case "IPADDR":
		v = IPADDR(buf.Bytes())
	case "TAG_INT":
		if length != 4 {
			return nil, ERR_ATT_TYPE
		}
		binary.Read(buf, binary.BigEndian, &v)
	case "TAG_STR":
		v = TAG_STR(buf.Bytes())
	case "HEXADECIMAL":
		v = HEXADECIMAL(buf.Bytes())
	default:
		return nil, ERR_ATT_TYPE
	}
	return v, nil
}

//根据指定内容生成属性值
//数字类型只接受int或uint32
//字符类型只接受string或[]byte
func NewAttributeValue(attrType string, i interface{}) (AttributeValue, error) {
	attrType = stringfix(attrType)
	v, err := newAttributeValue(attrType)
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
		return nil, ERR_ATT_TYPE
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
		return nil, ERR_ATT_TYPE
	case "IPADDR":
		tmp, ok := i.(string)
		if ok {
			v = IPADDR(tmp)
			return v, nil
		}
		tmp1, ok1 := i.([]byte)
		if ok1 {
			v = IPADDR(tmp1)
			return v, nil
		}
		return nil, ERR_ATT_TYPE
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
		return nil, ERR_ATT_TYPE
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
		return nil, ERR_ATT_TYPE
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
		return nil, ERR_ATT_TYPE
	default:
		return nil, ERR_ATT_TYPE
	}
}

//判断是否是有效的属性值类型字符串
func IsValidAttributeType(s string) bool {
	if s == "INTEGER" || s == "STRING" || s == "IPADDR" || s == "TAG_INT" || s == "TAG_STR" || s == "HEXADECIMAL" {
		return true
	}
	return false
}
