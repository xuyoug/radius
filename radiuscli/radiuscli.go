package radiuscli

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/xuyoug/radius"
	"math/rand"
	"net"
	"strconv"
	"strings"
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
//一个客户端只能有一个id_geter
type id_geter struct {
	id      radius.Id
	id_sync sync.Mutex
}

//
func (s *id_geter) getId() radius.Id {
	s.id_sync.Lock()
	if s.id == radius.Id(255) {
		s.id = radius.Id(0)
	} else {
		s.id++
	}
	n := s.id
	s.id_sync.Unlock()
	return n
}

//newidgeter
func newidgeter() *id_geter {
	idgter := new(id_geter)
	return idgter
}

//
//一个客户端可以有多个heaper
type heaper struct { //顺序执行获取，不考虑锁
	h [256]bool
	s sync.Mutex
}

//
func newheaper() *heaper {
	h := new(heaper)
	for i, _ := range h.h {
		h.h[i] = true
	}
	return h
}

//
func (h *heaper) setused(j int) error {
	if !h.h[j] {
		return errors.New("It's a used id")
	}
	h.s.Lock()
	h.h[j] = false
	h.s.Unlock()
	return nil
}

//
func (h *heaper) setfree(j int) error {
	if h.h[j] {
		return errors.New("It's a free id")
	}
	h.s.Lock()
	h.h[j] = true
	h.s.Unlock()
	return nil
}

//
func (h *heaper) isvalid(j int) bool {
	h.s.Lock()
	tmp := h.h[j]
	h.s.Unlock()
	return tmp
}

//keeper
type keeper struct {
	state        bool
	lastworktime time.Time
	lport        int
	svraddr      *net.UDPAddr
	conn         *net.UDPConn
	h            *heaper
	c_pre        chan *radius.Radius
	c_recive     [256]chan *radius.Radius
	c_close      chan bool
	C_err        chan error
}

//newkeeper
//一个keeper一个本地端口
func newkeeper(svrip *net.UDPAddr, C_err_in chan error) (*keeper, error) {

	k := new(keeper)
	k.svraddr = svrip
	con, err := net.DialUDP("udp4", nil, svrip)
	if err != nil {
		return nil, err
	}
	k.conn = con
	tmp_s := strings.Split(con.LocalAddr().String(), ":")
	k.lport, _ = strconv.Atoi(tmp_s[1])
	k.h = newheaper()
	k.c_pre = make(chan *radius.Radius, 256) //
	for i := 0; i < 256; i++ {               //
		k.c_recive[i] = make(chan *radius.Radius)
	}
	k.c_close = make(chan bool) //
	k.C_err = C_err_in          //
	k.state = true
	go k.keeperrecive()
	go k.keepersend()
	fmt.Println("create a new keeper", k.lport)
	return k, nil
}

//close
func (k *keeper) close() {
	fmt.Println("close keeper", k.lport)
	k.state = false
	//
	for {
		if len(k.c_pre) == 0 {
			close(k.c_pre)
			break
		}
	}
	close(k.c_pre)
	//
	for i := 0; i < 256; i++ {
		close(k.c_recive[i])
	}
	//
	close(k.c_close)
	//
	k.conn.Close()
}

//keepersend
func (k *keeper) keepersend() {
	for {
		select {
		case r := <-k.c_pre:
			_, err := k.conn.Write(r.Bytes())
			if err != nil {
				k.C_err <- err
			}
			k.lastworktime = time.Now()
		case <-k.c_close:
			fmt.Print("close keeper and keeper_sender out")
			break
		}
	}
}

//keeperrecive
func (k *keeper) keeperrecive() {
	for {
		data := make([]byte, 4096)
		n, udpaddr_in, err := k.conn.ReadFromUDP(data)
		if err != nil && k.state {
			k.C_err <- err
		}

		if udpaddr_in.String() != k.svraddr.String() {
			err := errors.New("WARNING:recived data from :" + udpaddr_in.String())
			k.C_err <- err
			break //recive the target server's data
		}
		r := radius.NewRadius()
		r.Read(bytes.NewBuffer(data[0:n]))
		id := int(r.Id)
		if k.h.isvalid(id) {
			err := errors.New("Recived response but aleady timeout: port:" + strconv.Itoa(k.lport) + " Id:" + strconv.Itoa(id))
			k.C_err <- err
			continue
		}

		k.c_recive[id] <- r
		if !k.state {
			break
		}
	}
}

//send
func (k *keeper) send(r *radius.Radius, timeout int) (*radius.Radius, time.Duration) {
	id := int(r.Id)
	t_d, _ := time.ParseDuration(strconv.Itoa(timeout) + "s")
	k.h.setused(id)
	k.c_pre <- r
	t1 := time.Now()
	select {
	case rr := <-k.c_recive[id]:
		t := time.Since(t1)
		k.h.setfree(id)
		return rr, t
	case <-time.After(t_d):
		k.h.setfree(id)
		err := errors.New("Recive response timeout: port:" + strconv.Itoa(k.lport) + " Id:" + strconv.Itoa(id))

		k.C_err <- err
		return nil, t_d
	}
	panic("panic in keeper send method")
}

//Sender
type RadiusSender struct {
	svrip   *net.UDPAddr
	idgtr   *id_geter
	keepers map[int]*keeper
	secret  string
	timeout int
	C_err   chan error
}

//NewSender
func NewSender(dstip string, port int, secret string, timeout int) (*RadiusSender, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", dstip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	fmt.Println(udpaddr)
	sdr := new(RadiusSender)
	sdr.svrip = udpaddr
	sdr.secret = secret
	sdr.idgtr = newidgeter()
	sdr.keepers = make(map[int]*keeper)
	sdr.C_err = make(chan error, 1024)
	sdr.timeout = timeout

	//
	go sdr.callbackkeeper()
	return sdr, nil

}

//Close
func (rs *RadiusSender) Close() {
	for i, v := range rs.keepers {
		delete(rs.keepers, i)
		v.close()
	}
	close(rs.C_err)
}

//callbackkeeper
func (rs *RadiusSender) callbackkeeper() {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			for i, v := range rs.keepers {
				if time.Since(v.lastworktime) > time.Second*10 { //
					delete(rs.keepers, i)
					v.close()
				}
			}
		}
	}
}

//newid
func (rs *RadiusSender) newId() radius.Id {
	return rs.idgtr.getId()
}

//newkeeper
func (rs *RadiusSender) newkeeper() int {
	k, err := newkeeper(rs.svrip, rs.C_err)
	if err != nil {
		panic(err.Error())
	}
	if _, ok := rs.keepers[k.lport]; ok {
		panic("panic in creat new keeper but it is here")
	}
	rs.keepers[k.lport] = k
	return k.lport
}

//getvalidkeeper
func (rs *RadiusSender) getvalidkeeper(id radius.Id) *keeper {
	i := int(id)
	if len(rs.keepers) == 0 { //
		kid := rs.newkeeper()
		return rs.keepers[kid]
	}
	for _, k := range rs.keepers { //
		if k.state && k.h.isvalid(i) {
			return k
		}
	}
	kid := rs.newkeeper() //
	return rs.keepers[kid]
}

//Send
func (rs *RadiusSender) Send(r *radius.Radius) (*radius.Radius, time.Duration, error) {
	r.SetAuthenticator(rs.secret)
	r.SetLength()
	r.Id = rs.newId()
	k := rs.getvalidkeeper(r.Id) //
	rr, t := k.send(r, rs.timeout)

	if rr == nil {
		return nil, t, errors.New("Timeout")
	}
	if !rr.IsResponseValid(r.Authenticator, rs.secret) { //验证authenticator
		return nil, t, errors.New("Recived but Authenticator Error")
	}
	return rr, t, nil
}

//
