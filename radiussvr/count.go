package radiussvr

import (
	"net"
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

//fmtgoroutine
func (c *RadiusListener) addfmtgoroutine() {
	c.lsr_sync_f.Lock()
	c.fmtgoroutine++
	c.lsr_sync_f.Unlock()
}

//fmtgoroutine
func (c *RadiusListener) endfmtgoroutine() {
	c.lsr_sync_f.Lock()
	c.fmtgoroutine--
	c.lsr_sync_f.Unlock()
}

//fmtgoroutine
func (c *RadiusListener) Count_fmtgoroutine() int {
	c.lsr_sync_f.RLock()
	n := c.fmtgoroutine
	c.lsr_sync_f.RUnlock()
	return n
}

//
func (c *RadiusListener) add_received(ip_in net.IP) {
	ip := ip_in.String()
	c.lsr_sync_r.Lock()
	c.cnt_received++ //更新总计数器
	//更新node计数
	_, ok := c.nodesreceived[ip]
	if ok {
		c.nodesreceived[ip]++
	}
	c.nodesreceived[ip] = 1
	c.lsr_sync_r.Unlock()
}

//
func (c *RadiusListener) add_replyed(ip_in net.IP) {
	ip := ip_in.String()
	c.lsr_sync_r.Lock()
	c.cnt_replyed++ //更新总计数器
	//更新node计数
	_, ok := c.nodesreplyed[ip]
	if !ok {
		c.nodesreplyed[ip] = 1
	}
	c.nodesreplyed[ip]++
	c.lsr_sync_r.Unlock()
}

//Add_wrong添加一个错误计数，并发送到错误队列
func (c *RadiusListener) Add_wrong(ip_in net.IP, err error) {
	ip := ip_in.String()
	c.lsr_sync_w.Lock()
	c.cnt_wrong++ //更新总计数器
	//更新node计数
	_, ok := c.nodeswrong[ip]
	if !ok {
		c.nodeswrong[ip] = make(map[error]int)
	}
	_, ok = c.nodeswrong[ip][err]
	if !ok {
		c.nodeswrong[ip][err] = 1
		return
	}
	c.nodeswrong[ip][err]++
	c.lsr_sync_w.Unlock()
	//
	c.sendErr(NodeErr{ip_in, err})
}

//
func (c *RadiusListener) Count_Discard() int {
	c.lsr_sync_r.RLock()
	c.lsr_sync_s.RLock()
	discarded := c.cnt_received - c.cnt_replyed
	c.lsr_sync_r.RUnlock()
	c.lsr_sync_s.RUnlock()
	return discarded
}

//
func (c *RadiusListener) Count_Received() int {
	c.lsr_sync_r.RLock()
	n := c.cnt_received
	c.lsr_sync_r.RUnlock()
	return n
}

//
func (c *RadiusListener) Count_Replyed() int {
	c.lsr_sync_s.RLock()
	n := c.cnt_replyed
	c.lsr_sync_s.RUnlock()
	return n
}

//
func (c *RadiusListener) Count_Wrong() int {
	c.lsr_sync_w.RLock()
	n := c.cnt_wrong
	c.lsr_sync_w.RUnlock()
	return n
}

//
func (c *RadiusListener) Count_NodeDiscard(ip_in net.IP) int {
	ip := ip_in.String()
	c.lsr_sync_r.RLock()
	c.lsr_sync_s.RLock()
	discarded := c.nodesreceived[ip] - c.nodesreplyed[ip]
	c.lsr_sync_r.RUnlock()
	c.lsr_sync_s.RUnlock()
	return discarded
}

//
func (c *RadiusListener) Count_NodeReceived(ip_in net.IP) int {
	ip := ip_in.String()
	c.lsr_sync_r.RLock()
	n := c.nodesreceived[ip]
	c.lsr_sync_r.RUnlock()
	return n
}

//
func (c *RadiusListener) Count_NodeReplyed(ip_in net.IP) int {
	ip := ip_in.String()
	c.lsr_sync_s.RLock()
	n := c.nodesreplyed[ip]
	c.lsr_sync_s.RUnlock()
	return n
}

//
func (c *RadiusListener) Count_NodeWrong(ip_in net.IP) map[error]int {
	ip := ip_in.String()
	c.lsr_sync_w.RLock()
	n := c.nodeswrong[ip]
	c.lsr_sync_w.RUnlock()
	return n
}

//
func (c *RadiusListener) Count_NodeWrongDesignated(ip_in net.IP, err error) int {
	ip := ip_in.String()
	c.lsr_sync_w.RLock()
	n, ok := c.nodeswrong[ip][err]
	if !ok {
		return 0
	}
	c.lsr_sync_w.RUnlock()
	return n
}
