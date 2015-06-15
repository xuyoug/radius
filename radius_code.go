package radius

import (
	"bytes"
	"strconv"
)

//定义radius的CODE 一个字节
type Code uint8

//定义常见的radiuscode
const (
	CodeErr Code = iota
	CodeAccessRequest
	CodeAccessAccept
	CodeAccessReject
	CodeAccountingRequest
	CodeAccountingRespons
	CodeAccessChallenge  Code = 11
	CodeStatusServer     Code = 12
	CodeStatusClient     Code = 13
	CodeDisconnectReq    Code = 40
	CodeDisconnectAck    Code = 41
	CodeDisconnectNak    Code = 42
	CodeChangeFiltersReq Code = 43
	CodeChangeFiltersAck Code = 44
	CodeChangeFiltersNak Code = 45
	//CodeReserved         Code = 255   //保留，不应该使用
)

//methods of Radius_Code
func (c Code) String() string {
	switch c {
	case CodeAccessRequest:
		return "AccessRequest(1)"
	case CodeAccessAccept:
		return "AccessAccept(2)"
	case CodeAccessReject:
		return "AccessReject(3)"
	case CodeAccountingRequest:
		return "AccountingRequest(4)"
	case CodeAccountingRespons:
		return "AccountingRespons(5)"
	case CodeAccessChallenge:
		return "AccessChallenge(11)"
	case CodeStatusServer:
		return "StatusServer(12)"
	case CodeStatusClient:
		return "StatusClient(13)"
	case CodeDisconnectReq:
		return "DisconnectReq(40)"
	case CodeDisconnectAck:
		return "DisconnectAck(41)"
	case CodeDisconnectNak:
		return "DisconnectNak(42)"
	case CodeChangeFiltersReq:
		return "ChangeFiltersReq(43)"
	case CodeChangeFiltersAck:
		return "ChangeFiltersAck(44)"
	case CodeChangeFiltersNak:
		return "ChangeFiltersNak(45)"
	default:
		return "NotSupportedRadius:Code(" + strconv.Itoa(int(c)) + ")"
	}
}

//返回Code的int表示
func (c Code) Int() int {
	if c.IsSupported() {
		return int(c)
	}
	return 0
}

//从buffer填充Code
func (c *Code) read(buf *bytes.Buffer) error {
	b, err := buf.ReadByte()
	if err != nil {
		return err
	}
	i := Code(b)

	*c = i
	return nil
}

//将Code写入buffer
func (c Code) write(buf *bytes.Buffer) {
	buf.WriteByte(byte(c))
}

//Judge判断响应报文的Code
func (r Code) ack(judge bool) Code {
	switch r {
	case CodeAccessRequest:
		if judge {
			return CodeAccessAccept
		}
		return CodeAccessReject
	case CodeAccountingRequest:
		return CodeAccountingRespons
	case CodeDisconnectReq:
		if judge {
			return CodeDisconnectAck
		}
		return CodeDisconnectNak
	case CodeChangeFiltersReq:
		if judge {
			return CodeChangeFiltersAck
		}
		return CodeChangeFiltersNak
	}
	return CodeErr
}

//判断是否是支持的Code
func (r Code) IsSupported() bool {
	if r == CodeAccessRequest || r == CodeAccessAccept || r == CodeAccessReject ||
		r == CodeAccountingRequest || r == CodeAccountingRespons ||
		r == CodeAccessChallenge ||
		r == CodeStatusServer || r == CodeStatusClient ||
		r == CodeDisconnectReq || r == CodeDisconnectAck || r == CodeDisconnectNak ||
		r == CodeChangeFiltersReq || r == CodeChangeFiltersAck || r == CodeChangeFiltersNak {
		return true
	}
	return false
}

//判断是否是请求报文
func (r Code) IsRequest() bool {
	if r == CodeAccessRequest || r == CodeAccountingRequest {
		return true
	}
	return false
}

//判断是否是响应报文
func (r Code) IsRespons() bool {
	if r == CodeAccessAccept || r == CodeAccessReject || r == CodeAccountingRespons {
		return true
	}
	return false
}

//判断是否是COA请求报文
func (r Code) IsCOARequest() bool {
	if r == CodeDisconnectReq || r == CodeChangeFiltersReq {
		return true
	}
	return false
}

//判断是否是COA响应报文
func (r Code) IsCOARespons() bool {
	if r == CodeDisconnectAck || r == CodeDisconnectNak || r == CodeChangeFiltersAck || r == CodeChangeFiltersNak {
		return true
	}
	return false
}
