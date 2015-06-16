package radiussvr

import (
	//"net"
	"sync/atomic"
	"time"
)

//AliveTime返回radiuslistener的生存时间
func (c *RadiusListener) AliveTime() time.Duration {
	return time.Since(c.startTime)
}

//StratTime返回radiuslistener的启动时间
func (c *RadiusListener) StratTime() time.Time {
	return c.startTime
}

//sendErr向lsnr的错误队列发送一条错误
func (c *RadiusListener) sendErr(err error) {
	c.C_err <- err
}

//
func (c *RadiusListener) add_received() {

	atomic.AddInt64(&c.cnt_received, 1) //更新总计数器
	//更新node计数
}

//
func (c *RadiusListener) add_replyed() {

	atomic.AddInt64(&c.cnt_replyed, 1) //更新总计数器
	//更新node计数

}

//Add_wrong添加一个错误计数，并发送到错误队列
func (c *RadiusListener) Add_wrong() {

	atomic.AddInt64(&c.cnt_wrong, 1) //更新总计数器

}

//
func (c *RadiusListener) Count_Discard() int64 {

	discarded := c.cnt_received - c.cnt_replyed

	return discarded
}

//
func (c *RadiusListener) Count_Received() int64 {
	n := c.cnt_received
	return n
}

//
func (c *RadiusListener) Count_Replyed() int64 {
	n := c.cnt_replyed
	return n
}

//
func (c *RadiusListener) Count_Wrong() int64 {
	n := c.cnt_wrong
	return n
}

// //
// func (c *RadiusListener) Count_NodeDiscard(ip_in interface{}) int {
// 	var ip string
// 	var v1, v2 int
// 	switch ip_in.(type) {
// 	case string:
// 		ip = ip_in.(string)
// 	case net.IP:
// 		ip = ip_in.(net.IP).String()
// 	}

// 	v1, ok1 := c.nodesreceived[ip]
// 	if !ok1 {
// 		v1 = 0
// 	}
// 	v2, ok2 := c.nodesreplyed[ip]
// 	if !ok2 {
// 		v2 = 0
// 	}
// 	discarded := v1 - v2

// 	return discarded
// }

// //
// func (c *RadiusListener) Count_NodeReceived(ip_in interface{}) int {
// 	var ip string
// 	switch ip_in.(type) {
// 	case string:
// 		ip = ip_in.(string)
// 	case net.IP:
// 		ip = ip_in.(net.IP).String()
// 	}
// 	c.lsr_sync_r.RLock()
// 	n, ok := c.nodesreceived[ip]
// 	if !ok {
// 		n = 0
// 	}
// 	c.lsr_sync_r.RUnlock()
// 	return n
// }

// //
// func (c *RadiusListener) Count_NodeReplyed(ip_in interface{}) int {
// 	var ip string
// 	switch ip_in.(type) {
// 	case string:
// 		ip = ip_in.(string)
// 	case net.IP:
// 		ip = ip_in.(net.IP).String()
// 	}
// 	c.lsr_sync_s.RLock()
// 	n, ok := c.nodesreplyed[ip]
// 	if !ok {
// 		n = 0
// 	}
// 	c.lsr_sync_s.RUnlock()
// 	return n
// }

// //
// func (c *RadiusListener) Count_NodeWrong(ip_in interface{}) map[error]int {
// 	var ip string
// 	switch ip_in.(type) {
// 	case string:
// 		ip = ip_in.(string)
// 	case net.IP:
// 		ip = ip_in.(net.IP).String()
// 	}
// 	c.lsr_sync_w.RLock()
// 	n, ok := c.nodeswrong[ip]
// 	if !ok {
// 		n = nil
// 	}
// 	c.lsr_sync_w.RUnlock()
// 	return n
// }

// //
// func (c *RadiusListener) Count_NodeWrongDesignated(ip_in interface{}, err error) int {
// 	var ip string
// 	switch ip_in.(type) {
// 	case string:
// 		ip = ip_in.(string)
// 	case net.IP:
// 		ip = ip_in.(net.IP).String()
// 	}
// 	c.lsr_sync_w.RLock()
// 	n, ok := c.nodeswrong[ip][err]
// 	if !ok {
// 		return 0
// 	}
// 	c.lsr_sync_w.RUnlock()
// 	return n
// }
