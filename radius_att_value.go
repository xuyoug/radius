package radius

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//"reflect"
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
	return strconv.Itoa(int(v))
}
func (v TAG_STR) String() string {
	return string(v)
}
func (v HEXADECIMAL) String() string {
	return fmt.Sprintf("%v", []byte(v))
}

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

//获取值
func (v INTEGER) WriteBuff(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v STRING) WriteBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v IPADDR) WriteBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v TAG_INT) WriteBuff(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v TAG_STR) WriteBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v HEXADECIMAL) WriteBuff(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}

func FillValueFromBuff(v *AttributeValue, buf *bytes.Buffer) {
	fmt.Println(v, *v, buf.Bytes())
	fmt.Println((*v).(INTEGER))
	switch (*v).Typestring() {
	case "INTEGER":
		var tmp int
		binary.Read(buf, binary.BigEndian, &tmp)
		*v = INTEGER(tmp)
	case "STRING":
		*v = STRING(buf.Bytes())
	case "IPADDR":
		*v = IPADDR(buf.Bytes())
	case "TAG_INT":
		var tmp int
		binary.Read(buf, binary.BigEndian, &tmp)
		*v = INTEGER(tmp)
	case "TAG_STR":
		*v = TAG_STR(buf.Bytes())
	case "HEXADECIMAL":
		*v = HEXADECIMAL(buf.Bytes())
	default:
		fmt.Println("here")
	}
	fmt.Println(v, *v, buf.Bytes())
}

//定时属性的值类型
type AttributeValue interface {
	//fillValue(buf *bytes.Buffer)
	WriteBuff(buf *bytes.Buffer)
	String() string
	Value() interface{}
	Len() uint8
	Typestring() string
}

func NewAttributeValue(attrType string) (AttributeValue, error) {
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

func NewAttributeValueFromBuff(attrType string, buf *bytes.Buffer) (AttributeValue, error) {
	attrType = stringfix(attrType)
	switch attrType {
	case "INTEGER":
		var v INTEGER
		binary.Read(buf, binary.BigEndian, &v)
		return v, nil
	case "STRING":
		var v STRING
		v = STRING(buf.Bytes())
		return v, nil
	case "IPADDR":
		var v IPADDR
		v = IPADDR(buf.Bytes())
		return v, nil
	case "TAG_INT":
		var v TAG_INT
		binary.Read(buf, binary.BigEndian, &v)
		return v, nil
	case "TAG_STR":
		var v TAG_STR
		v = TAG_STR(buf.Bytes())
		return v, nil
	case "HEXADECIMAL":
		var v HEXADECIMAL
		v = HEXADECIMAL(buf.Bytes())
		return v, nil
	default:
		return nil, ERR_ATT_TYPE
	}
}

//
func isValidAttributeType(s string) bool {
	if s == "INTEGER" || s == "STRING" || s == "IPADDR" || s == "TAG_INT" || s == "TAG_STR" || s == "HEXADECIMAL" {
		return true
	}
	return false
}
