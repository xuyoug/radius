package radiussvr

import (
	"errors"
	"net"
)

var (
	Err_InvaildNode  = errors.New("Invalid Node")
	Err_SecretWrong  = errors.New("Secret Wrong")
	Err_FmtError     = errors.New("Fmt Radius Error")
	Err_ReplyTimeout = errors.New("Reply Timeout")
	Err_Drop_SrcChan = errors.New("Drop Inside SrcChan Because Timeout")
	Err_CanotReply   = errors.New("Can not reply this radius")
)

type NodeErr struct {
	IP  net.IP
	Err error
}

func (ne *NodeErr) String() string {
	return ne.IP.String() + ":" + ne.Err.Error()
}

func (ne NodeErr) Error() string {
	return ne.IP.String() + ":" + ne.Err.Error()
}

func NewNodeErr(ip net.IP, err error) *NodeErr {
	ne := new(NodeErr)
	ne.IP = ip
	ne.Err = err
	return ne
}

type SvrErr struct {
	Discript string
	Err      error
}

func (se *SvrErr) String() string {
	return se.Discript + ":" + se.Err.Error()
}

func (se SvrErr) Error() string {
	return se.Discript + ":" + se.Err.Error()
}
