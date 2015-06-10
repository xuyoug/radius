package radiuscli

import (
	"bytes"
	"crypto/md5"
	"errors"
	"github.com/xuyoug/radius"
	"math/rand"
	"net"
	"sync"
	"time"
)

//radius客户端的实现封装

var cli_source rand.Source
var cli_rand *rand.Rand

//
func init() {
	cli_source = rand.NewSource(int64(time.Now().Nanosecond()))
	cli_rand = rand.New(cli_source)
}

//计算随机数
func getrand(i int) int {
	return cli_rand.Intn(i)
}

//
func newAuthenticator() radius.R_Authenticator {
	bs := make([]byte, 16)
	for i := 0; i < 16; i++ {
		bs = append(bs, byte(getrand(255)))
	}
	return radius.R_Authenticator(bs)
}

//
func setAuthenticator(r *radius.Radius, secret string) radius.R_Authenticator {
	if r.R_Code == radius.CodeAccountingRequest {
		buf := bytes.NewBuffer([]byte{})
		r.R_Authenticator = radius.R_Authenticator([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		r.WriteToBuff(buf)
		buf.Write(secret)
		m := md5.Sum(buf.Bytes())
		r.R_Authenticator = radius.R_Authenticator(m[:])
		return r.R_Authenticator
	}
	r.R_Authenticator = newAuthenticator()
	return r.R_Authenticator
}

//
//一个客户端只能有一个id_geter
type id_geter struct {
	id      radius.R_Id
	id_sync sync.Mutex
}

//
func (s *id_geter) getId() radius.R_Id {
	s.id_sync.Lock()
	if s.id == radius.R_Id(255) {
		s.id = 0
	} else {
		s.id++
	}
	n := s.id
	s.id_sync.Unlock()
	return n
}

//
//一个客户端可以有多个heaper
type heaper struct { //顺序执行获取，不考虑锁
	h [256]bool
}

//
func newheaper() *heaper {
	h := new(heaper)
	for i, _ := range h.h {
		h.h[i] = true
	}
}

//
func (h *heaper) pop() (radius.R_Id, error) {
	for i, v := range h.h {
		if v {
			return radius.R_Id(i), nil
		}
	}
	return radius.R_Id(0), errors.New("No valid id in this heaper")
}

//
func (h *heaper) push(i radius.R_Id) {
	j := int(i)
	if !h.h[j] {
		h.h[j] = true
	}
}

//
func (h *heaper) isvalid(i radius.R_Id) bool {
	j := int(i)
	return h.h[j]
}

//keeper
type keeper struct {
	raddr    *net.UDPAddr
	lport    int
	conn     *net.UDPConn
	h        *heaper
	c_pre    chan radius.Radius
	c_send   chan radius.Radius
	c_recive [256]chan radius.Radius
}

//newkeeper

//close

//Sender
type RadiusSender struct {
	id_geter
	keepers []keeper
	secret  string
	c_send  chan radius.Radius
}

//NewSender
func NewSender(dstip string, port int, secret string, timeout time.Duration) (*RadiusSender, error) {

}

//Close
func (rs *RadiusSender) Close() {

}

//Send
func (rs *RadiusSender) Send(r *radius.Radius) (*radius.Radius, error) {

}

//
