package radius

//定义radius的基础数据结构
type Radius struct {
	Code
	Id
	Length
	Authenticator
	AttributeList
}

//methods of Radus
func (r *Radius) String() string {
	return r.Code.String() + "\n" +
		r.Id.String() + "\n" +
		r.Length.String() + "\n" +
		r.Authenticator.String() + "\n" +
		r.AttributeList.String()
}

//初始化一个radius
func NewRadius() *Radius {
	r := new(Radius)
	r.AttributeList.attributes = make([]Attribute, 0)
	return r
}

//初始化一个radius
func NewRadiusI(in interface{}) *Radius {
	r := NewRadius()

	switch in.(type) {
	case Code:
		r.Code = in.(Code)
	case int:
		r.Code = Code(in.(int))
	default:
		return nil
	}
	if r.Code == CodeAccessRequest {
		r.Authenticator = newAuthenticator()
	}
	return r
}

//Ack创建响应radius
func (r *Radius) Ack(judge bool) *Radius {
	response := r.ack(judge)
	if response == CodeErr {
		return nil
	}
	rr := NewRadiusI(response)
	rr.Authenticator = r.Authenticator
	rr.Id = r.Id
	return rr
}
