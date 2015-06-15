package radius

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

//
//定义radius中各结构的方法
//

//checkLengthWithBuff判断buf长度和radius Length是否相等
func (r *Radius) checkLengthWithBuff(buf *bytes.Buffer) bool {
	l := R_Length(buf.Len())
	if r.R_Length == l {
		return true
	}
	return false
}

//从buf填充radius Length
func (r *R_Length) readFromBuff(buf *bytes.Buffer) error {
	var b1, b2 byte
	var err1, err2 error
	b1, err1 = buf.ReadByte()
	b2, err2 = buf.ReadByte()
	if err1 != nil || err2 != nil {
		return ERR_LEN_INVALID
	}
	l := R_Length(b1<<8) + R_Length(b2)
	if l.IsValidLenth() && buf.Len()+4 >= int(l) { //不允许buf长度小于radius长度，但是大于可以
		*r = l
		return nil
	}
	return ERR_LEN_INVALID
}

//methods of Radus
func (r *Radius) String() string {
	return r.Code.String() + "\n" +
		r.R_Id.String() + "\n" +
		r.R_Length.String() + "\n" +
		r.R_Authenticator.String() + "\n" +
		r.AttributeList.String()
}

//ReadFromBuffer从buf填充radius结构
func (r *Radius) ReadFromBuffer(buf *bytes.Buffer) error {
	err := r.Code.readFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on Code")
	}

	err = r.R_Id.readFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on Id")
	}

	err = r.R_Length.readFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on Length")
	}

	err = r.R_Authenticator.readFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on Authenticator")
	}
	for {
		v, err := readAttribute(buf)
		if isEOF(err) {
			break
		}
		if err != nil {
			return err
		}
		r.AttributeList.AddAttr(v)
	}
	if r.GetLength() != r.R_Length {
		return ERR_OTHER
	}
	return nil
}

//WriteToBuff将radius结构字节化写入buf
func (r *Radius) WriteToBuff(buf *bytes.Buffer) {
	buf.WriteByte(byte(r.Code))
	buf.WriteByte(byte(r.R_Id))
	binary.Write(buf, binary.BigEndian, r.R_Length)
	buf.Write([]byte(r.R_Authenticator))
	for _, v := range r.AttributeList.attributes {
		v.writeBuffer(buf)
	}
}

//AddAttribute添加radius属性id-value值对
func (r *Radius) AddAttr(attv interface{}, attid ...interface{}) error {
	Attid, err := NewAttributeId(attid)
	if err != nil {
		return err
	}
	typ := Attid.ValueTypestring()
	Attv, err1 := NewAttributeValue(typ, attv)
	if err1 != nil {
		return err1
	}

	r.AttributeList.AddAttr(Attribute{Attid, Attv})
	return nil
}

//AddAttribute添加radius属性id-value值对
func (r *Radius) AddAttrS(attv_s string, attid ...interface{}) error {
	Attid, err := NewAttributeId(attid)
	if err != nil {
		return err
	}
	typ := Attid.ValueTypestring()
	Attv, err1 := NewAttributeValueS(typ, attv_s)
	if err1 != nil {
		return err1
	}

	r.AttributeList.AddAttr(Attribute{Attid, Attv})
	return nil
}
