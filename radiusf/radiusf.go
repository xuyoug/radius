package radiusf

import (
	"bytes"
	"errors"
	//"fmt"
	//"github.com/xuyoug/radius"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

//def the err
var (
	ERR_FMT = errors.New("fmt package err")
)

//def the attribute struct
//it do not show for outside
type attr struct {
	sid uint8
	v   uint32
	vid uint8
}

//def the method to get attribute's name
func (a *attr) Name() string {
	if v, ok := map_f[*a]; ok {
		return v.Name
	}
	return ""
}

//def the attribute's value type
func (a *attr) Type() string {
	if v, ok := map_f[*a]; ok {
		return v.Type
	}
	return ""
}

//def radiusf struct
type Radiusf struct {
	Code  uint8
	attrs map[string]interface{}
}

//format the radiusf type from buffer
//return nill and err when get a terrible error
//return formatable part and err when get other error
func ReadRadiusf(buf *bytes.Buffer) (*Radiusf, error) {
	var b [4]byte
	var err error
	var l, l1 int
	var vname string
	for i := 0; i < 4; i++ {
		b[i], err = buf.ReadByte()
		if err != nil {
			return nil, ERR_FMT
		}
	}
	r := new(Radiusf)
	r.attrs = make(map[string]interface{})
	r.Code = uint8(b[0])
	l = int(b[2])<<8 + int(b[3])
	buf.Next(16)
	buf = bytes.NewBuffer(buf.Next(l - 20))
	for {
		b[0], err = buf.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return r, ERR_FMT
		}
		a := new(attr)
		a.sid = uint8(b[0])

		b[0], err = buf.ReadByte()
		if err != nil {
			return r, ERR_FMT
		}
		l = int(b[0])

		if a.sid != 26 {
			vname = a.Name()
			bs := buf.Next(l - 2)
			if vname == "" {
				continue
			}

			switch a.Type() {
			case "STRING":
				r.attrs[vname] = string(bs)
			case "INTEGER":
				var i int
				binary.Read(bytes.NewBuffer(bs), binary.BigEndian, &i)
				r.attrs[vname] = i
			case "IPADDR":
				r.attrs[vname] = bs
			}
		} else {
			for i := 0; i < 4; i++ {
				b[i], err = buf.ReadByte()
				if err != nil {
					return r, ERR_FMT
				}
			}
			a.v = uint32(b[0])<<24 + uint32(b[1])<<16 + uint32(b[2])<<8 + uint32(b[3])
			if !isvalidv(a.v) {
				continue
			}
			b[0], err = buf.ReadByte()
			if err != nil {
				return r, ERR_FMT
			}
			a.vid = uint8(b[0])
			b[1], err = buf.ReadByte()
			if err != nil {
				return r, ERR_FMT
			}
			l1 = int(b[1])
			if l1 != l-2 {
				return r, ERR_FMT
			}

			vname = a.Name()
			bs := buf.Next(l1 - 6)
			if vname == "" {
				continue
			}

			switch a.Type() {
			case "STRING":
				r.attrs[vname] = string(bs)
			case "INTEGER":
				var i int
				binary.Read(bytes.NewBuffer(bs), binary.BigEndian, &i)
				r.attrs[vname] = i
			case "IPADDR":
				r.attrs[vname] = bs
			}
		}
	}
	return r, nil
}

//get attribute by name and return the value as it's type
//return nil when not found
func (r *Radiusf) Get(attname string) interface{} {
	attname = strings.ToUpper(attname)
	if v, ok := r.attrs[attname]; ok {
		return v
	}
	return nil
}

//get attribute by name and format the value as a formatable string
func (r *Radiusf) GetS(attname string) string {
	attname = strings.ToUpper(attname)
	if v, ok := r.attrs[attname]; ok {
		switch v.(type) {
		case int:
			return strconv.Itoa(v.(int))
		case string:
			return v.(string)
		case []byte:
			return string(v.([]byte))
		}
	}
	return ""
}
