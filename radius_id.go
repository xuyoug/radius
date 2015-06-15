package radius

import (
	"bytes"
	"strconv"
)

//定义radius的ID
//一个字节
type Id uint8

//methods of R_Id
func (i Id) String() string {
	return "Id(" + strconv.Itoa(int(i)) + ")"
}

//从buffer填充Id
func (id *Id) read(buf *bytes.Buffer) error {
	b, err := buf.ReadByte()
	if err != nil {
		return err
	}
	i := Id(b)
	*id = i
	return nil
}

//将Id写入buffer
func (id Id) write(buf *bytes.Buffer) {
	err := buf.WriteByte(byte(id))
}
