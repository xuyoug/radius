package radius

import (
	"encoding/binary"
	"strconv"
)

//定义radius的Len 两个字节
type Length uint16

//定义radiusLength的最大最小值
const (
	Length_MIN Length = 20
	Length_MAX Length = 4096
)

//methods of R_Length
func (l Length) String() string {
	return "Length(" + strconv.Itoa(int(l)) + ")"
}

//判断是否是有效的radius长度
func (r Length) IsValidLenth() bool {
	if r >= R_Length_MIN && r <= R_Length_MAX {
		return true
	}
	return false
}

//从buffer填充Length
func (r *Length) read(buf *bytes.Buffer) error {
	var b1, b2 byte
	var err1, err2 error
	b1, err1 = buf.ReadByte()
	b2, err2 = buf.ReadByte()
	if err1 != nil || err2 != nil {
		return ERR_RADIUS_FMT
	}
	l := R_Length(b1<<8) + R_Length(b2)
	if l.IsValidLenth() && buf.Len()+4 >= int(l) { //不允许buf长度小于radius长度，但是大于可以
		*r = l
		return nil
	}
	return ERR_LEN_INVALID
}

//将Length写入buffer
func (r Length) write(buf *bytes.Buffer) {
	binary.Write(buf, binary.BigEndian, r)
}
