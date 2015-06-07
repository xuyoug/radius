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

//methods of Radius_Code
func (i R_Code) String() (s string) {
	switch i {
	case CodeAccessRequest:
		return "CodeAccessRequest(1)"
	case CodeAccessAccept:
		return "CodeAccessAccept(2)"
	case CodeAccessReject:
		return "CodeAccessReject(3)"
	case CodeAccountingRequest:
		return "CodeAccountingRequest(4)"
	case CodeAccountingRespons:
		return "CodeAccountingRespons(5)"
	case CodeAccessChallenge:
		return "CodeAccessChallenge(11)"
	case CodeStatusServer:
		return "CodeStatusServer(12)"
	case CodeStatusClient:
		return "CodeStatusClient(13)"
	case CodeReserved:
		return "CodeReserved(255)"
	}
	return ERR_CODE_WRONG.Error() + ":(" + strconv.Itoa(int(i)) + ")"
}

//
func (r *R_Code) readFromBuff(buf *bytes.Buffer) error {
	b, err := buf.ReadByte()
	if err != nil {
		return ERR_RADIUS_FMT
	}
	i := R_Code(b)
	if i < 6 || (i >= 11 && i <= 13) || i == 255 {
		*r = i
		return nil
	}
	return ERR_CODE_WRONG
}

//
func (r R_Code) Judge(judge bool) (R_Code, error) {
	switch r {
	case CodeAccessRequest:
		if judge {
			return CodeAccessAccept, nil
		}
		return CodeAccessReject, nil
	case CodeAccountingRequest:
		return CodeAccountingRespons, nil
	}
	return CodeAccessReject, ERR_NOTSUPPORT
}

//
func (r R_Code) IsSupported() bool {
	if r == CodeAccessRequest || r == CodeAccessAccept || r == CodeAccessReject || r == CodeAccountingRequest || r == CodeAccountingRespons {
		return true
	}
	return false
}

//
func (r R_Code) IsRequest() bool {
	if r == CodeAccessRequest || r == CodeAccountingRequest {
		return true
	}
	return false
}

//
func (r R_Code) IsRespons() bool {
	if r == CodeAccessAccept || r == CodeAccessReject || r == CodeAccountingRespons {
		return true
	}
	return false
}

//methods of R_Id
func (i R_Id) String() string {
	return fmt.Sprintf("Id(%d)", i)
}

//
func (r *R_Id) readFromBuff(buf *bytes.Buffer) error {
	b, err := buf.ReadByte()
	if err != nil {
		return ERR_RADIUS_FMT
	}
	i := R_Id(b)
	*r = i
	return nil
}

//methods of R_Length
func (l R_Length) String() string {
	return fmt.Sprintf("Length(%d)", l)
}

//
func (r R_Length) isValidLenth() bool {
	if r >= radiusLength_MIN || r <= radiusLength_MAX {
		return true
	}
	return false
}

//
func (r *Radius) checkLengthWithBuff(buf *bytes.Buffer) bool {
	l := R_Length(buf.Len())
	if r.R_Length == l {
		return true
	}
	return false
}

//
func (r *R_Length) readFromBuff(buf *bytes.Buffer) error {
	var b1, b2 byte
	var err1, err2 error
	b1, err1 = buf.ReadByte()
	b2, err2 = buf.ReadByte()
	if err1 != nil || err2 != nil {
		return ERR_LEN_INVALID
	}
	l := R_Length(b1<<8) + R_Length(b2)
	if l.isValidLenth() {
		*r = l
		return nil
	}
	return ERR_LEN_INVALID
}

//methods of R_Authenticator
func (a R_Authenticator) String() string {
	return fmt.Sprintf("Authenticator %v", []byte(a))
}

//
func (r *R_Authenticator) readFromBuff(buf *bytes.Buffer) error {
	b := buf.Next(Radius_Authenticator_LEN)
	*r = b
	return nil
}

//

//methosd of Attributes maps of Id or Name

//methods of Radus
func (r *Radius) String() string {
	return r.R_Code.String() + "\n" +
		r.R_Id.String() + "\n" +
		r.R_Length.String() + "\n" +
		r.R_Authenticator.String() + "\n" +
		r.AttributeList.String()
}

//
func (r *Radius) ReadFromBuffer(buf *bytes.Buffer) error {
	err := r.R_Code.readFromBuff(buf)
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
	return nil
}

//
func (r *Radius) WriteToBuff(buf *bytes.Buffer) {
	buf.WriteByte(byte(r.R_Code))
	buf.WriteByte(byte(r.R_Id))
	binary.Write(buf, binary.BigEndian, r.R_Length)
	buf.Write([]byte(r.R_Authenticator))
	for _, v := range r.AttributeList.attributes {
		v.writeBuffer(buf)
	}
}

//
func (r *Radius) GetLength() R_Length {
	var l R_Length
	l = 20
	for _, v := range r.AttributeList.attributes {
		switch v.AttributeId.(type) {
		case AttId:
			l += R_Length(v.AttributeValue.ValueLen() + 2)
		case AttIdV:
			if v.AttributeId.(AttIdV).VendorTypestring() == "IETF" {
				l += R_Length(v.AttributeValue.ValueLen() + 8)
			}
			if v.AttributeId.(AttIdV).VendorTypestring() == "TYPE4" {
				l += R_Length(v.AttributeValue.ValueLen() + 10)
			}
		}
	}
	return l
}

//
func (r *Radius) setLength() {
	r.R_Length = r.GetLength()
}

//
func (r *Radius) setAuthenticator() {

}

//
func (r *Radius) Finished() {
	r.setLength()
	r.setAuthenticator()
}
