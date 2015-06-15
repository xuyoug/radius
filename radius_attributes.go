package radius

import (
	"strconv"
)

//AttrList定义radius报文中的属性列表
type AttributeList struct {
	attributes []Attribute
}

//AddAttr在radius属性列表中添加一个明确的属性
func (a *AttributeList) AddAttr(r *Attribute) {
	a.attributes = append(a.attributes, *r)
}

//AddAttrEx在radius属性列表中创建并添加一个明确的属性
func (a *AttributeList) AddAttrEx(in ...interface{}) error {
	v, err := NewAttribute(in...)
	if err != nil {
		return err
	}
	a.AddAttr(v)
	return nil
}

//DelAttr在radius属性列表中根据attributeid删除属性
//返回删除的属性数量
func (a *AttributeList) DelAttr(r AttributeId) int {
	var n int
	vs := make([]Attribute, 0)
	for _, v := range a.attributes {
		if v.AttributeId != r {
			vs = append(vs, v)
		} else {
			n++
		}
	}
	a.attributes = vs
	return n
}

//GetAttrsNum获取AttributeList的属性数量
func (a *AttributeList) GetAttrsNum() int {
	return len(a.attributes)
}

//GetAttrNum获取AttributeList中指定attributeid的属性数量
func (a *AttributeList) GetAttrNum(r AttributeId) int {
	var n int
	for _, v := range a.attributes {
		if v.AttributeId == r {
			n++
		}
	}
	return n
}

//ContainsAttr获取AttributeList中指定attributeid的属性数量
func (a *AttributeList) ContainsAttr(r AttributeId) bool {
	for _, v := range a.attributes {
		if v.AttributeId == r {
			return true
		}
	}
	return false
}

//GetAttrs获取所有属性
func (a *AttributeList) GetAttrs() []Attribute {
	return a.attributes
}

//GetAttrValues获取指定属性的值的列表
func (a *AttributeList) GetAttrValues(r AttributeId) []AttributeValue {
	list := make([]AttributeValue, 0)
	for _, v := range a.attributes {
		if v.AttributeId == r {
			list = append(list, v.AttributeValue)
		}
	}
	return list
}

//GetAttrValue1获取第一个指定属性，若指定属性不存在，则返回错误
func (a *AttributeList) GetAttrValue1(r AttributeId) AttributeValue {
	for _, v := range a.attributes {
		if v.AttributeId == r {
			return v.AttributeValue
		}
	}
	return nil
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
		a.AddAttr(&Attribute{r, vs})
	}
}

//String方法打印属性列表自身
func (a *AttributeList) String() string {
	var s string
	if len(a.attributes) == 0 {
		return "No Attributes"
	}
	s += "Attributes:"
	for i, v := range a.attributes {
		s += "\n"
		s += strconv.Itoa(i+1) + ": " + v.String()
	}
	return s
}
