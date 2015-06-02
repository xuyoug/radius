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

//定时属性的值类型
type AttributeValue interface {
	//fillValue(buf *bytes.Buffer)
	writetobuf(buf *bytes.Buffer)
	String() string
	Len() uint8
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

//获取值
func (v *INTEGER) fillValue(buf *bytes.Buffer) {
	binary.Read(buf, binary.BigEndian, &v)
}
func (v *STRING) fillValue(buf *bytes.Buffer) {
	*v = STRING(buf.Bytes())
}
func (v *IPADDR) fillValue(buf *bytes.Buffer) {
	*v = IPADDR(buf.Bytes())
}
func (v *TAG_INT) fillValue(buf *bytes.Buffer) {
	binary.Read(buf, binary.BigEndian, &v)
}
func (v *TAG_STR) fillValue(buf *bytes.Buffer) {
	*v = TAG_STR(buf.Bytes())
}
func (v *HEXADECIMAL) fillValue(buf *bytes.Buffer) {
	*v = HEXADECIMAL(buf.Bytes())
}

//获取值
func (v INTEGER) writetobuf(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v STRING) writetobuf(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v IPADDR) writetobuf(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v TAG_INT) writetobuf(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v TAG_STR) writetobuf(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v HEXADECIMAL) writetobuf(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}

//
func isValidAttributeType(s string) bool {
	if s == "INTEGER" || s == "STRING" || s == "IPADDR" || s == "TAG_INT" || s == "TAG_STR" || s == "HEXADECIMAL" {
		return true
	}
	return false
}
