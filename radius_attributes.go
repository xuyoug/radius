package radius

//AttrList定义radius报文中的属性列表
type AttributeList struct {
	attributes []Attribute
}

//AddAttr在radius属性列表中添加一个明确的属性
func (a *AttributeList) AddAttr(r Attribute) {
	a.attributes = append(a.attributes, r)
}

//GetAttrsNum获取AttributeList的属性数量
func (a *AttributeList) GetAttrsNum() int {
	return len(a.attributes)
}

//GetAttrs获取所有属性
func (a *AttributeList) GetAttrs() []Attribute {
	return a.attributes
}

//GetAttr获取指定属性的值的列表
func (a *AttributeList) GetAttrValues(r AttributeId) []AttributeValue {
	list := make([]AttributeValue, 0)
	for _, v := range a.attributes {
		if v.AttributeId == r {
			list = append(list, v.AttributeValue)
		}
	}
	return list
}

//SetAttr设置指定属性的值
//若不存在，则添加一个该属性
func (a *AttributeList) SetAttr(r AttributeId, vs AttributeValue) {
	ischanged := false
	for _, v := range a.attributes {
		if v.AttributeId == r {
			v.AttributeValue = vs
			ischanged = true
		}
	}
	if !ischanged {
		a.AddAttr(Attribute{r, vs})
	}
}

//GetAttrFist获取第一个指定属性，若指定属性不存在，则返回错误
func (a *AttributeList) GetAttrValueFirst(r AttributeId) (AttributeValue, error) {
	for _, v := range a.attributes {
		if v.AttributeId == r {
			return v.AttributeValue, nil
		}
	}
	return nil, ERR_ATT_NO
}

//String方法打印属性列表自身
func (a *AttributeList) String() string {
	var s string
	s += "Attributes:"
	for i, v := range a.attributes {
		s += "\n"
		s += strconv.Itoa(i+1) + ": " + v.String()
	}
	return s
}
