package radius

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

//定义attributevalue的类型
type INTEGER uint32
type STRING string
type IPADDR net.IP
type TAG_INT uint32
type TAG_STR string
type HEXADECIMAL []byte

//Uint32返回IPADDR类型的uint32形式
//用于索引含义
func (i IPADDR) Uint32() uint32 {
	return uint32(i[0])<<24 + uint32(i[1])<<16 + uint32(i[2])<<8 + uint32(i[3])
}

//String格式化属性值输出
func (v INTEGER) String() string {
	return strconv.Itoa(int(v))
}
func (v STRING) String() string {
	return string(v)
}
func (v IPADDR) String() string {
	return net.IP(v).String()
}
func (v TAG_INT) String() string {
	return "TAG(" + strconv.Itoa(int(v>>24)) + ")" + strconv.Itoa(int(v))
}
func (v TAG_STR) String() string {
	return string(v)
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
	return 4
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
	return net.IP(v)
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
func (v INTEGER) Type() string {
	return "INTEGER"
}
func (v STRING) Type() string {
	return "STRING"
}
func (v IPADDR) Type() string {
	return "IPADDR"
}
func (v TAG_INT) Type() string {
	return "TAG_INT"
}
func (v TAG_STR) Type() string {
	return "TAG_STR"
}
func (v HEXADECIMAL) Type() string {
	return "HEXADECIMAL"
}

//writeBuff写值入buffer
func (v INTEGER) write(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v STRING) write(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v IPADDR) write(buf *bytes.Buffer) {
	buf.Write([]byte(net.IP(v).To4()))
}
func (v TAG_INT) write(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, v)
}
func (v TAG_STR) write(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}
func (v HEXADECIMAL) write(buf *bytes.Buffer) {
	buf.Write([]byte(v))
}

//AttributeValue定时属性值类型接口
type AttributeValue interface {
	write(buf *bytes.Buffer)
	String() string
	Value() interface{}
	ValueLen() int
	Type() string
}

//newAttributeValue获取空属性值
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
		return nil, ERR_ATT_VTYPE
	}
}

//readAttributeValueFromBuff从buff获取属性值
func readAttributeValue(attrType string, length int, buf_in *bytes.Buffer) (AttributeValue, error) {
	v, err := newAttributeValue(attrType)
	if err != nil {
		return nil, err
	}
	bs := buf_in.Next(length)
	buf := bytes.NewBuffer(bs)
	var tmp uint32
	switch attrType {
	case "INTEGER":
		if length != 4 {
			return nil, ERR_VALUE_TYPE
		}
		binary.Read(buf, binary.BigEndian, &tmp)
		v = INTEGER(tmp)
	case "STRING":
		v = STRING(bs)
	case "IPADDR":
		if length != 4 {
			return nil, ERR_VALUE_TYPE
		}
		v = IPADDR(bs)
	case "TAG_INT":
		if length != 4 {
			return nil, ERR_VALUE_TYPE
		}
		binary.Read(buf, binary.BigEndian, &tmp)
		v = TAG_INT(tmp)
	case "TAG_STR":
		v = TAG_STR(bs)
	case "HEXADECIMAL":
		v = HEXADECIMAL(bs)
	default:
		return nil, ERR_ATT_VTYPE
	}
	return v, nil
}

//NewAttributeValue根据指定内容生成属性值
//数字类型只接受int或uint32或可格式化为数字的字符串
//字符类型只接受string或[]byte
//该方法主要由外部使用
func NewAttributeValue(attrType string, i interface{}) (AttributeValue, error) {
	v, err := newAttributeValue(attrType)
	if err != nil {
		return nil, err
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
		tmp2, ok2 := i.(string)
		if ok2 {
			num, err2 := strconv.Atoi(tmp2)
			if err2 == nil {
				v = INTEGER(num)
				return v, nil
			}
		}
		return nil, ERR_VALUE_TYPE
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
		return nil, ERR_VALUE_TYPE
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
		if ok1 && len(tmp1) == 4 {
			v = IPADDR(tmp1)
			return v, nil
		}
		return nil, ERR_VALUE_TYPE
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
		tmp2, ok2 := i.(string)
		if ok2 {
			num, err2 := strconv.Atoi(tmp2)
			if err2 == nil {
				v = TAG_INT(num)
				return v, nil
			}
		}
		return nil, ERR_VALUE_TYPE
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
		return nil, ERR_VALUE_TYPE
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
		return nil, ERR_VALUE_TYPE
	default:
		return nil, ERR_ATT_VTYPE
	}
	//return nil, ERR_ATT_VTYPE
}

//NewAttributeValueS根据指定字符串生成属性值
//该方法主要由外部使用
func NewAttributeValueS(attrType, s string) (AttributeValue, error) {
	attrType = stringfix(attrType)
	v, err := newAttributeValue(attrType)
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
		return nil, ERR_ATT_VTYPE
	case "STRING":
		v = STRING(s)
		return v, nil
	case "IPADDR":
		p := net.ParseIP(s)
		if p != nil {
			v = IPADDR([]byte(p.To4()))
		} else {
			if len(s) == 4 {
				v = IPADDR(s)
			} else {
				return nil, ERR_ATT_VTYPE
			}
		}
		return v, nil
	case "TAG_INT":
		i, err1 := strconv.Atoi(s)
		if err1 == nil {
			v = INTEGER(i)
			return v, nil
		}
		return nil, ERR_ATT_VTYPE
	case "TAG_STR":
		v = TAG_STR(s)
		return v, nil
	case "HEXADECIMAL":
		v = HEXADECIMAL(s)
		return v, nil
	default:
		return nil, ERR_ATT_VTYPE
	}
	return nil, ERR_ATT_VTYPE
}

//IsValidAttributeType判断是否是有效的属性值类型字符串
func IsValidAttributeValueType(s string) bool {
	if s == "INTEGER" || s == "STRING" || s == "IPADDR" || s == "TAG_INT" || s == "TAG_STR" || s == "HEXADECIMAL" {
		return true
	}
	return false
}
