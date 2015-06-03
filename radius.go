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
	//"fmt"
	//"net"
)

// func GetRadius(conn *net.UDPConn) (*Radius, error) {
// 	var inbytes [4096]byte
// 	n, addr, err := conn.ReadFromUDP(inbytes[0:])
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(n, addr)
// 	fmt.Println(inbytes[0:n])

// 	r := NewRadius()
// 	r.FillFromBuf(bytes.NewBuffer(inbytes[0:n]))

// 	fmt.Println(r)

// 	//conn.WriteToUDP(inbytes[0:n], addr)

// 	return r, err
// }

//
func NewRadius() *Radius {
	r := new(Radius)
	r.AttributeList.list_name = make(map[AttributeId][]AttributeValue, 0)
	r.AttributeList.list_order = make([]AttributeId, 0)
	return r
}

func (r *Radius) FillFromBuf(buf *bytes.Buffer) error {

	err := r.getCodeFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_code")
	}

	err = r.getIdFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_Id")
	}

	err = r.getLenFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_Length")
	}

	err = r.getAuthenticatorFromBuff(buf)
	if err != nil {
		return errors.New("Format wrong on R_Length")
	}

	err = r.getAttsFromBuff(buf)
	if err != nil {
		return err
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
				v.Write(buf)
				buf.WriteByte(byte(vv.Len() + 2))
				vv.WriteBuff(buf)
			case AttIdV:
				buf.WriteByte(byte(ATTID_VENDOR_SPECIFIC))
				buf.WriteByte(byte(vv.Len() + 8))
				v.Write(buf)
				buf.WriteByte(byte(vv.Len() + 2))
				vv.WriteBuff(buf)
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
			case AttIdV:
				if v.(AttIdV).VendorId.Typestring() == "IETF" {
					l += R_Length(vv.Len() + 8)
				}
				if v.(AttIdV).VendorId.Typestring() == "TYPE4" {
					l += R_Length(vv.Len() + 8)
				}
			}
		}
	}
	return l
}
