package radius

//
//定义radius报文的结构化方法和序列化方法
//

//radius.go
//radius_att_def.go
//radius_att_meaning.go
//radius_att_value.go
//radius_att_vendor.go
//radius_com.go
//radius_def.go
//radius_errors.go
//radius_method.go

//定义radius的结构化方法
//定义radius的处理方法

//NewRadius初始化一个radius对象
func NewRadius() *Radius {
	r := new(Radius)
	r.R_Authenticator = make([]byte, 16)
	r.AttributeList.attributes = make([]Attribute, 0)
	return r
}

//NewRadiusI初始化一个确定类型的radius对象
func NewRadiusI(i interface{}) (*Radius, error) {
	r := NewRadius()
	switch i.(type) {
	case int:
		r.R_Code = R_Code(i.(int))
	case R_Code:
		r.R_Code = i.(R_Code)
	default:
		return nil, ERR_NOTSUPPORT
	}
	if !r.R_Code.IsSupported() {
		return nil, ERR_NOTSUPPORT
	}
	return r, nil
}

//Finish完成radius对象的长度封装
func (r *Radius) SetLength() {
	r.R_Length = r.GetLength()
}

//GetLength获取radius结构字节化后的长度
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
