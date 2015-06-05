package radiussvr

import (
	"bytes"
	"github.com/xuyoug/radius"
	"net"
	"strconv"
	"sync"
	"time"
)

//radius服务器的实现封装

//
type SrcRadius struct {
	SrcAddr    *net.UDPAddr
	Secret     string
	ReciveTime time.Time
	Radius     *radius.Radius
	lisenter   *RadiusListener
}

//
type ReplyRadius struct {
	DstAddr    *net.UDPAddr
	Secret     string
	ReciveTime time.Time
	Radius     *radius.Radius
	lisenter   *RadiusListener
}

//
func (sr *SrcRadius) Reply() *ReplyRadius {
	rr := new(ReplyRadius)
	rr.DstAddr = sr.SrcAddr
	rr.Secret = sr.Secret
	rr.ReciveTime = sr.ReciveTime
	rr.Radius = radius.NewRadius()
	rr.Radius.R_Id = sr.Radius.R_Id
	//计算Authenticator

	//计算
	rr.lisenter = sr.lisenter
	return rr
}

//
func (dr *ReplyRadius) Send() {
	dr.lisenter.c_send <- dr
}

//
type RadiusListener struct {
	conn       *net.UDPConn
	udpAddr    *net.UDPAddr
	secretlist *SecretList
	Recived    int
	Send       int
	c_recive   chan SrcRadius
	c_send     chan SrcRadius
	c_err      chan error
	startTime  time.Time
	timeout    time.Duration
	lsr_sync   sync.RWMutex
}

//
func (c *RadiusListener) run(cache int) error {
	c.c_recive = make(chan SrcRadius, cache)
	c.c_send = make(chan ReplyRadius, cache)
	c.c_err = make(chan error, cache)
	c.startTime = time.Now()
	con, err := net.ListenMulticastUDP("udp", nil, c.udpAddr)
	if err != nil {
		return err
	}
	c.conn = con
	go c.getSrcRadius()
	go c.replyRadius()
}

//
func (c *RadiusListener) getSrcRadius() {
	var bs [4096]byte
	var b_num int
	var udpAddr *net.UDPAddr
	var ip net.IP
	var secret string
	var err error
	for {
		b_num, udpAddr, err := c.conn.ReadFromUDP(bs[0:])
		//
		c.lsr_sync.Lock()
		c.c_recive++
		c.lsr_sync.Unlock()
		//
		if b_num > 4096 {
			err = radius.ERR_LEN_INVALID
		}
		if err != nil {
			c.c_err <- err
			return
		}
		//
		buf := bytes.NewBuffer(bs[0:b_num])
		r := radius.NewRadius()
		err = r.ReadFromBuffer(buf)
		if err != nil {
			c.c_err <- err
			return
		}
		//
		ip = udpAddr.IP
		secret = c.secretlist.GetSecret(ip)
		//
		src_r := new(SrcRadius)
		src_r.SrcAddr = udpAddr
		src_r.ReciveTime = time.Now()
		src_r.Secret = secret
		src_r.Radius = r
		src_r.lisenter = c
		//
		select {
		case c.c_recive <- src_r:
			return
		case <-time.After(time.Second):
			c.c_err <- ERR_DROP_TO
		}
	}
}

//
func (c *RadiusListener) replyRadius() {

}

//
func (c *RadiusListener) AliveTime() time.Duration {
	return time.Since(c.startTime)
}

//
func (c *RadiusListener) Count_Discard() int {
	c.lsr_sync.RLock()
	discarded := c.Recived - c.Send
	c.lsr_sync.RUnlock()
	return discarded
}

//
func (c *RadiusListener) Count_Recived() int {
	c.lsr_sync.RLock()
	n := c.Recived
	c.lsr_sync.RUnlock()
	return n
}

//
func (c *RadiusListener) Count_Send() int {
	c.lsr_sync.RLock()
	n := c.Send
	c.lsr_sync.RUnlock()
	return n
}

//
func RadiusServer(localip string, port int, secret_list *SecretList, timeout time.Duration, cache int) (*RadiusListener, error) {
	addr := localip + ":" + strconv.Itoa(port)
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	lsr := new(RadiusListener)
	lsr.udpAddr = udpAddr
	lsr.secretlist = secret_list
	lsr.timeout = timeout

	err = lsr.run(cache)
	if err != nil {
		return nil, err
	}

	return lsr, nil

}
