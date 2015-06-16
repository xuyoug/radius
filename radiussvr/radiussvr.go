package radiussvr

import (
	"bytes"
	//"fmt"
	"github.com/xuyoug/radius"
	"net"
	"strconv"
	"time"
)

//radius服务器的实现封装

//SrcRadius定义收到的radius源报文结构
type SrcRadius struct {
	SrcAddr    *net.UDPAddr
	Secret     string
	ReciveTime time.Time
	Radius     *radius.Radius
	lisenter   *RadiusListener
}

//ReplyRadius定义响应radius报文结构
type ReplyRadius struct {
	DstAddr    *net.UDPAddr
	Secret     string
	ReciveTime time.Time
	Radius     *radius.Radius
	lisenter   *RadiusListener
}

//Reply由源报文生成基础响应报文
//不带任何属性
func (sr *SrcRadius) Reply(judge bool) (*ReplyRadius, error) {
	rr := new(ReplyRadius)
	rr.DstAddr = sr.SrcAddr
	rr.Secret = sr.Secret
	rr.ReciveTime = sr.ReciveTime
	rr.lisenter = sr.lisenter
	rr.Radius = sr.Radius.Ack(judge)
	//
	if rr.Radius == nil {
		rr.lisenter.Add_wrong()
		return nil, Err_CanotReply
	}
	//
	return rr, nil //然后交由外部处理
}

//Send发送响应报文
func (dr *ReplyRadius) Send() {
	dr.lisenter.c_send <- dr
}

//定义radiuslistener的结构
type RadiusListener struct {
	conn         *net.UDPConn
	udpAddr      *net.UDPAddr
	secretlist   *SecretList
	cnt_received int64
	cnt_replyed  int64
	cnt_wrong    int64
	C_recive     chan *SrcRadius
	c_or         chan *original_radius
	c_send       chan *ReplyRadius
	C_err        chan error
	startTime    time.Time
	timeout      time.Duration
}

//定义从网卡获取的原始为格式化的radius信息
type original_radius struct {
	udpAddr *net.UDPAddr
	buf     *bytes.Buffer
}

//run启动RadiusListener
func (c *RadiusListener) run(cache int) error {
	c.C_recive = make(chan *SrcRadius, cache)
	c.c_or = make(chan *original_radius, cache)
	c.c_send = make(chan *ReplyRadius, cache)
	c.C_err = make(chan error, cache)
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

//getSrcOriginalbytes从网卡获取的原始为格式化的radius信息
//并将其存入缓存chan
func (c *RadiusListener) getSrcOriginalbytes() {
	for {
		var bs [4096]byte
		var b_num int
		var udpAddr *net.UDPAddr
		var err error
		b_num, udpAddr, err = c.conn.ReadFromUDP(bs[0:])
		c.add_received()
		if err != nil {
			c.Add_wrong()
			return
		}
		if b_num > 4096 || b_num < 20 { //对于长度非法的自己忽略
			err = radius.ERR_LEN_INVALID
			//fmt.Println("length error")
			c.Add_wrong()
			return
		}
		or := new(original_radius)
		or.buf = bytes.NewBuffer(bs[0:b_num])
		or.udpAddr = udpAddr
		c.c_or <- or
	}
}

//fmtRadius从缓存chan获取原始radius报文
//一旦获取数据则启动一个新goroutine对其格式化
func (c *RadiusListener) fmtRadius() {
	for {
		select {
		case or := <-c.c_or:
			go c.decoderadius(or)
		}
	}
}

//decoderadius格式化原始radius报文
func (c *RadiusListener) decoderadius(or *original_radius) {
	var ip net.IP
	var secret string
	var err error
	//添加一个格式化计数
	//c.addfmtgoroutine()
	src_r := new(SrcRadius)

	ip = or.udpAddr.IP                           //获取IP
	secret = c.secretlist.GetSecret(ip.String()) //根据IP获取密钥
	src_r.Secret = secret
	//

	src_r.SrcAddr = or.udpAddr
	src_r.ReciveTime = time.Now()
	src_r.lisenter = c
	src_r.Radius, err = radius.ReadFromBuffer(or.buf)

	if err != nil {
		c.Add_wrong()
		c.C_err <- err
		//c.endfmtgoroutine()
		return
	}
	//如果是计费请求报文，验证authenticator
	if !src_r.CheckAuthenticator() {
		c.Add_wrong()
		//c.endfmtgoroutine()
		return
	}
	//
	select {
	case c.C_recive <- src_r:
		//c.endfmtgoroutine()
		return
	case <-time.After(time.Second):
		//c.endfmtgoroutine()
		c.Add_wrong()
	}
}

//作为radiuslistener的发送goroutine
//将发送队列内的响应报文发送出去
func (c *RadiusListener) replyRadius() {
	var err error
	for {
		select {
		case rr := <-c.c_send:
			rr.Radius.SetAuthenticator(rr.Secret)
			rr.Radius.SetLength()
			_, err = c.conn.WriteToUDP(rr.Radius.Bytes(), rr.DstAddr)
			if err != nil {
				rr.lisenter.Add_wrong()
			}
			c.add_replyed()
		}
	}
}

//RadiusServer生成一个radiuslistener并启动
//传入参数：启动端口、密钥管理器、超时、队列深度
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
