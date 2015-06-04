package radiussvr

import (
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
func (sr *SrcRadius) ReplyRadius() *ReplyRadius {
	rr := new(ReplyRadius)
	rr.DstAddr = sr.SrcAddr
	rr.Secret = sr.Secret
	rr.ReciveTime = sr.ReciveTime
	rr.Radius = radius.NewRadius()
	rr.lisenter = sr.lisenter
	return rr
}

//
func (sr *ReplyRadius) Send() {
	sr.lisenter.c_send <- sr
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
	startTime  time.Time
	timeout    time.Duration
	lsr_sync   sync.RWMutex
}

//
func (c *RadiusListener) run(cache int) error {
	c.C_recive = make(chan SrcRadius, cache)
	c.C_send = make(chan SrcRadius, cache)
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
	for {

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
