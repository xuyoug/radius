package radius

import (
	"bytes"
	"net"
)

//
//定义radius的方法
//
func ReadFromBuffer(buf *bytes.Buffer) (*Radius, error) {
	r := NewRadius()
	err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//
func (r *Radius) Read(buf *bytes.Buffer) error {
	err := r.Code.read(buf)
	if err != nil {
		return err
	}

	err = r.Id.read(buf)
	if err != nil {
		return err
	}

	err = r.Length.read(buf)
	if err != nil {
		return err
	}

	err = r.Authenticator.read(buf)
	if err != nil {
		return err
	}
	for {
		v, err := readAttribute(buf)
		if isEOF(err) {
			break
		}
		if err != nil {
			return err
		}
		r.AttributeList.AddAttr(&v)
	}

	return nil
}

//Write将radius结构字节化写入buf
func (r *Radius) Write(buf *bytes.Buffer) {
	r.SetLength()
	r.Code.write(buf)
	r.Id.write(buf)
	r.Length.write(buf)
	r.Authenticator.write(buf)
	for _, v := range r.AttributeList.attributes {
		v.write(buf)
	}
}

//Bytes将radius序列化为[]byte
func (r *Radius) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	r.Write(buf)
	return buf.Bytes()
}

//Send设置radius的authenticator和length
//然后将其发送到网络上
func (r *Radius) Send(c *net.UDPConn, secret string) error {
	r.SetAuthenticator(secret)
	r.SetLength()
	_, err := c.Write(r.Bytes())
	if err != nil {
		return err
	}
	return nil
}
