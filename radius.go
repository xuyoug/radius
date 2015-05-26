package radius

//
//定义radius的结构化方法和处理方法
//

//定义radius的结构化方法

//定义radius的处理方法

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func GetRadius(conn *net.UDPConn) (*Radius, error) {
	var inbytes [4096]byte
	n, addr, err := conn.ReadFromUDP(inbytes[0:])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(n, addr)
	fmt.Println(inbytes[0:n])

	r := NewRadius()
	r.FillFromBuf(bytes.NewBuffer(inbytes[0:n]))

	fmt.Println(r)

	//conn.WriteToUDP(inbytes[0:n], addr)

	return r, err
}

//
func NewRadius() *Radius {
	r := new(Radius)
	r.AttributeList.list_name = make(map[AttributeId][]AttributeValue, 0)
	r.AttributeList.list_order = make([]AttributeId, 0)
	return r
}

func (r *Radius) FillFromBuf(buf *bytes.Buffer) error {

	err := r.getRadiusCodeFbuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_code")
	}

	err = r.getRadiusIdFbuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_Id")
	}

	err = r.getRadiusLenFbuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_Length")
	}

	err = r.getRadiusAuthenticatorFbuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_Length")
	}

	err = r.getRadiusAttsFbuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_Attributes")
	}

	return nil
}

//
func (r *Radius) WriteToBuff() *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(byte(r.R_Code))
	buf.WriteByte(byte(r.R_Id))
	binary.Write(buf, binary.BigEndian, r.R_Length)
	buf.Write([]byte(r.R_Authenticator))
	for _, v := range r.AttributeList.list_order {
		for _, vv := range r.AttributeList.list_name[v] {
			switch v.(type) {
			case AttId:
				buf.WriteByte(byte(v.(AttId)))
				buf.WriteByte(byte(vv.Len() + 2))
				vv.writetobuf(buf)
			case AttVId:
				buf.WriteByte(byte(VENDOR_SPECIFIC))
				buf.WriteByte(byte(vv.Len() + 8))
				v.writevendor(buf)
				v.writeAtt(buf)
				buf.WriteByte(byte(vv.Len() + 2))
				vv.writetobuf(buf)
			case AttV4Id:
				buf.WriteByte(byte(VENDOR_SPECIFIC))
				buf.WriteByte(byte(vv.Len() + 10))
				v.writevendor(buf)
				v.writeAtt(buf)
				vv.writetobuf(buf)
			}
		}
	}
	return buf
}

//
func (r *Radius) GetLength() R_Length {
	var l R_Length
	l = 20
	for _, v := range r.AttributeList.list_order {
		for _, vv := range r.AttributeList.list_name[v] {
			switch v.(type) {
			case AttId:
				l += R_Length(vv.Len() + 2)
			case AttVId:
				l += R_Length(vv.Len() + 8)
			case AttV4Id:
				l += R_Length(vv.Len() + 10)
			}
		}
	}
	return l
}
