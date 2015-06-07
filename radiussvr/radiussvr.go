package radiussvr

import (
	"bytes"
	//"fmt"
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
	buf        *bytes.Buffer
}

//
type ReplyRadius struct {
	DstAddr    *net.UDPAddr
	Secret     string
	ReciveTime time.Time
	Radius     *radius.Radius
	lisenter   *RadiusListener
	buf        *bytes.Buffer
}

//
func (sr *SrcRadius) Reply(judge bool) (*ReplyRadius, error) {
	var err error
	rr := new(ReplyRadius)
	rr.DstAddr = sr.SrcAddr
	rr.Secret = sr.Secret
	rr.ReciveTime = sr.ReciveTime
	rr.lisenter = sr.lisenter
	rr.Radius = radius.NewRadius()
	rr.Radius.R_Id = sr.Radius.R_Id
	rr.Radius.R_Authenticator = sr.Radius.R_Authenticator
	rr.Radius.R_Code, err = sr.Radius.R_Code.Judge(judge)
	//
	if err != nil {
		rr.lisenter.Add_wrong(rr.DstAddr.IP, err)
		return nil, err
	}
	//
	rr.buf = bytes.NewBuffer([]byte{})

	return rr, nil //然后交由外部处理
}

//
func (rr *ReplyRadius) makebuf() {
	rr.Radius.R_Length = rr.Radius.GetLength()
	rr.Radius.WriteToBuff(rr.buf)
	//计算最新的authenticator
	rr.Radius.R_Authenticator = radius.R_Authenticator(rr.ReplyAuthenticator())

	//
	rr.buf = bytes.NewBuffer([]byte{})
	rr.Radius.WriteToBuff(rr.buf)
}

//
func (dr *ReplyRadius) Send() {
	dr.makebuf()
	dr.lisenter.c_send <- dr
}

//
type RadiusListener struct {
	conn          *net.UDPConn
	udpAddr       *net.UDPAddr
	secretlist    *SecretList
	cnt_received  int
	cnt_replyed   int
	cnt_wrong     int
	nodesreceived map[string]int
	nodesreplyed  map[string]int
	nodeswrong    map[string]map[error]int
	fmtgoroutine  int //标识当前有多少个协程在解radius报文
	C_recive      chan *SrcRadius
	c_or          chan *original_radius
	c_send        chan *ReplyRadius
	C_err         chan error
	startTime     time.Time
	timeout       time.Duration
	lsr_sync_r    sync.RWMutex
	lsr_sync_s    sync.RWMutex
	lsr_sync_w    sync.RWMutex
	lsr_sync_f    sync.RWMutex
}

//
type original_radius struct {
	udpAddr *net.UDPAddr
	buf     *bytes.Buffer
}

//
func (c *RadiusListener) run(cache int) error {
	c.C_recive = make(chan *SrcRadius, cache)
	c.c_or = make(chan *original_radius, cache)
	c.c_send = make(chan *ReplyRadius, cache)
	c.C_err = make(chan error, cache)
	c.nodesreceived = make(map[string]int)
	c.nodesreplyed = make(map[string]int)
	c.nodeswrong = make(map[string]map[error]int)
	c.startTime = time.Now()
	con, err := net.ListenUDP("udp", c.udpAddr)
	if err != nil {
		return err
	}
	c.conn = con
	go c.getSrcOriginalbytes()
	go c.fmtRadius()
	go c.replyRadius()
	return nil
}

//
func (c *RadiusListener) getSrcOriginalbytes() {
	for {
		var bs [4096]byte
		var b_num int
		var udpAddr *net.UDPAddr
		var err error
		b_num, udpAddr, err = c.conn.ReadFromUDP(bs[0:])
		c.add_received(udpAddr.IP)
		if err != nil {
			c.Add_wrong(udpAddr.IP, err)
			return
		}
		if b_num > 4096 || b_num < 20 {
			err = radius.ERR_LEN_INVALID
			c.Add_wrong(udpAddr.IP, err)
			return
		}
		or := new(original_radius)
		or.buf = bytes.NewBuffer(bs[0:b_num])
		or.udpAddr = udpAddr
		c.c_or <- or
	}
}

//
func (c *RadiusListener) fmtRadius() {
	for {
		select {
		case or := <-c.c_or:
			go c.decoderadius(or)
		}
	}
}

//
func (c *RadiusListener) decoderadius(or *original_radius) {
	var ip net.IP
	var secret string
	var err error
	//
	c.addfmtgoroutine()
	src_r := new(SrcRadius)

	ip = or.udpAddr.IP
	secret = c.secretlist.GetSecret(ip.String())
	src_r.Secret = secret
	//

	src_r.SrcAddr = or.udpAddr
	src_r.ReciveTime = time.Now()
	src_r.lisenter = c
	src_r.buf = or.buf
	src_r.Radius = radius.NewRadius()

	err = src_r.Radius.ReadFromBuffer(src_r.buf)
	if err != nil {
		c.Add_wrong(ip, err)
		c.endfmtgoroutine()
		return
	}
	//如果是计费请求报文，验证authenticator
	if !src_r.IsValidAuthenticator() {
		c.Add_wrong(ip, Err_SecretWrong)
		c.endfmtgoroutine()
		return
	}
	//
	select {
	case c.C_recive <- src_r:
		c.endfmtgoroutine()
		return
	case <-time.After(time.Second):
		c.endfmtgoroutine()
		c.Add_wrong(ip, Err_Drop_SrcChan)
	}
}

//
func (c *RadiusListener) replyRadius() {
	var err error
	for {
		select {
		case rr := <-c.c_send:
			_, err = rr.lisenter.conn.WriteToUDP(rr.buf.Bytes(), rr.DstAddr)
			if err != nil {
				rr.lisenter.Add_wrong(rr.DstAddr.IP, err)
			}
			c.add_replyed(rr.DstAddr.IP)
		}
	}
}

//
func RadiusServer(port int, secret_list *SecretList, timeout time.Duration, cache int) (*RadiusListener, error) {
	addr := ":" + strconv.Itoa(port)
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	//
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
