package radius

import (
	"bytes"
	"fmt"
	"io"
	//"strings"
	"encoding/binary"
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

func (r *Radius) getRadiusCodeFbuff(buf *bytes.Buffer) error {
	b, err := buf.ReadByte()
	if err != nil {
		return ERR_RADIUS_FMT
	}
	i := R_Code(b)
	if i < 6 || (i >= 11 && i <= 13) || i == 255 {
		r.R_Code = i
		return nil
	}
	return ERR_CODE_WRONG
}

//methods of R_Id

func (i R_Id) String() string {
	return fmt.Sprintf("Id(%d)", i)
}

func (r *Radius) getRadiusIdFbuff(buf *bytes.Buffer) error {
	b, err := buf.ReadByte()
	if err != nil {
		return ERR_RADIUS_FMT
	}
	i := R_Id(b)
	r.R_Id = i
	return nil
}

//methods of R_Length

func (l R_Length) String() string {
	return fmt.Sprintf("Length(%d)", l)
}

func (l R_Length) IsValidRLenth() bool {
	if l >= radiusLength_MIN && l <= radiusLength_MAX {
		return true
	}
	return false
}

func (r *Radius) checkRadiusLength(buf *bytes.Buffer) bool {
	l := R_Length(buf.Len())
	if r.R_Length == l {
		return true
	}
	return false
}

func (r *Radius) getRadiusLenFbuff(buf *bytes.Buffer) error {
	var b1, b2 byte
	var err1, err2 error
	b1, err1 = buf.ReadByte()
	b2, err2 = buf.ReadByte()
	if err1 != nil || err2 != nil {
		return ERR_LEN_INVALID
	}
	l := R_Length(b1<<8) + R_Length(b2)
	if l.IsValidRLenth() {
		r.R_Length = l
		return nil
	}
	return ERR_LEN_INVALID
}

//methods of R_Authenticator

func (a R_Authenticator) String() string {
	return fmt.Sprintf("Authenticator(%v)", []byte(a))
}

func (r *Radius) getRadiusAuthenticatorFbuff(buf *bytes.Buffer) error {
	b := buf.Next(Radius_Authenticator_LEN)
	if len(b) != Radius_Authenticator_LEN {
		return ERR_AUTHENTICATOR_INVALID
	}
	r.R_Authenticator = b
	return nil
}

//

//methosd of Attributes maps of Id or Name

//methods of Attributes
func (r *Radius) getRadiusAttsFbuff(buf *bytes.Buffer) error {

	for {
		var b1, b2, b3, b4 byte
		var err error
		var tmplen int
		b1, err = buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return ERR_ATT_FMT
		}
		b2, err = buf.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return ERR_ATT_FMT
		}
		tmplen = int(b2) - 2
		if int(b1) != 26 {
			vid := AttId(b1)
			vtype := vid.Typestring()

			v, err := NewAttributeValueFromBuff(vtype, bytes.NewBuffer(buf.Next(tmplen)))
			if err != nil {
				return err
			}
			r.AttributeList.AddAttr(vid, v)
		} else {
			var vid VendorId
			binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &vid)
			if !vid.IsvalidVendor() {
				return ERR_VENDOR_INVALID
			}
			if vid.Typestring() == "IETF" {
				b3, err = buf.ReadByte()
				if err != nil {
					if err == io.EOF {
						break
					}
					return ERR_ATT_FMT
				}
				vaid := AttV(b3)
				va := AttVId{vid, vaid}

				b4, err = buf.ReadByte()
				if err != nil {
					if err == io.EOF {
						break
					}
					return ERR_ATT_FMT
				}

				la := int(b4 - 2)
				v, err := NewAttributeValueFromBuff(va.Typestring(), bytes.NewBuffer(buf.Next(la)))
				if err != nil {
					return err
				}
				r.AttributeList.AddAttr(va, v)

			}
			if vid.Typestring() == "TYPE4" {
				var vaid AttV4
				binary.Read(bytes.NewBuffer(buf.Next(4)), binary.BigEndian, &vaid)
				va := AttV4Id{vid, vaid}

				v, err := NewAttributeValueFromBuff(va.Typestring(), bytes.NewBuffer(buf.Next(tmplen-8)))
				if err != nil {
					return err
				}
				r.AttributeList.AddAttr(va, v)
			}
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
