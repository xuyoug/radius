package radius

import (
	"bytes"
	"fmt"
	"io"
	//"strings"
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

func (r *R_Code) getCodeFromBuff(buf *bytes.Buffer) error {
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

//methods of R_Id

func (i R_Id) String() string {
	return fmt.Sprintf("Id(%d)", i)
}

func (r *R_Id) getIdFromBuff(buf *bytes.Buffer) error {
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

func (r R_Length) isValidLenth() bool {
	if r >= radiusLength_MIN || r <= radiusLength_MAX {
		return true
	}
	return false
}

func (r *Radius) checkRadiusLengthWithBuff(buf *bytes.Buffer) bool {
	l := R_Length(buf.Len())
	if r.R_Length == l {
		return true
	}
	return false
}

func (r *R_Length) getLenFromBuff(buf *bytes.Buffer) error {
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

func (r *R_Authenticator) getAuthenticatorFromBuff(buf *bytes.Buffer) error {
	b := buf.Next(Radius_Authenticator_LEN)
	if len(b) != Radius_Authenticator_LEN {
		return ERR_AUTHENTICATOR_INVALID
	}
	*r = b
	return nil
}

//

//methosd of Attributes maps of Id or Name

//methods of Attributes
func (r *Radius) getAttsFromBuff(buf *bytes.Buffer) error {
	//循环读取buff，直至结束或错误
	for {
		var blen byte
		var err error
		var tmplen int
		var tmplen_v int
		var aid AttId

		aid, err = readAttId(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return ERR_ATT_FMT
		}

		blen, err = buf.ReadByte()
		if err != nil {
			return ERR_ATT_FMT
		}

		tmplen = int(blen) - 2
		if aid != ATTID_VENDOR_SPECIFIC {
			vtype := aid.Typestring()
			v, err := NewAttributeValueFromBuff(vtype, tmplen, buf)
			if err != nil {
				return err
			}
			r.AttributeList.AddAttr(aid, v)
		} else {
			aidv, err1 := readAttIdV(buf)
			if err1 != nil {
				return err1
			}
			if aidv.VendorId.Typestring() == "IETF" {
				blen, err = buf.ReadByte()
				if err != nil {
					return ERR_ATT_FMT
				}
				tmplen_v = int(blen) - 2
				if tmplen_v != tmplen-6 {
					return ERR_ATT_FMT
				}
			}
			if aidv.VendorId.Typestring() == "TYPE4" {
				tmplen_v = tmplen - 8
			}
			v, err := NewAttributeValueFromBuff(aidv.Typestring(), tmplen_v, buf)
			if err != nil {
				return err
			}
			r.AttributeList.AddAttr(aidv, v)
		}
	}
	return nil
}

//methods of Radus
func (r *Radius) String() string {
	return r.R_Code.String() + "\n" +
		r.R_Id.String() + "\n" +
		r.R_Length.String() + "\n" +
		r.R_Authenticator.String() + "\n" +
		r.AttributeList.String()
}
