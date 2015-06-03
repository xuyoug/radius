package radius

//
//定义radius的结构化方法和处理方法
//

//定义radius的结构化方法

//定义radius的处理方法

import ()

// func GetRadius(conn *net.UDPConn) (*Radius, error) {
// 	var inbytes [4096]byte
// 	n, addr, err := conn.ReadFromUDP(inbytes[0:])
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	fmt.Println(n, addr)
// 	fmt.Println(inbytes[0:n])

// 	r := NewRadius()
// 	r.FillFromBuf(bytes.NewBuffer(inbytes[0:n]))

// 	fmt.Println(r)

// 	//conn.WriteToUDP(inbytes[0:n], addr)

// 	return r, err
// }

//
func NewRadius() *Radius {
	r := new(Radius)
	r.AttributeList.attributes = make([]Attribute, 0)
	return r
}
