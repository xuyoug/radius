package radius

import (
	"bytes"
)

//AttId定义标准属性Id
type AttId uint8

//String方法返回AttId其名字
func (a AttId) String() string {
	s, ok := list_attributestand_id[a]
	if ok {
		return s.Name
	}
	return "UNKNOWN_STAND_ATTRIBUTE(" + strconv.Itoa(int(a)) + ")"
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
func (a AttId) ValueType() string {
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

//write方法将AttId自己写入buffer
func (a AttId) write(buf *bytes.Buffer) error {
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
		return ATTID_ERR, err
	}
	return AttId(b), nil
}

//GetAttId提供根据名字返回AttId的方法
func GetAttId(s string) AttId {
	s = stringfix(s)
	a, ok := list_attributestand_name[s]
	if ok {
		return a
	}
	return ATTID_ERR
}
