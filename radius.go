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
	r.AttributeList.attributes = make([]Attribute, 0)
	return r
}
