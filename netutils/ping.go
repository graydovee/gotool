package netutils

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func Ping(address string, timeout time.Duration) *PingResult {
	hosts, err := net.LookupHost(address)
	if err != nil {
		return &PingResult{
			TargetAddr: address,
			Err:        err,
		}
	}
	var res *PingResult
	for _, host := range hosts {
		res = Pinger(host, timeout)
		if res.Err == nil {
			return res
		}
	}
	if res.Err == nil {
		res.Err = errors.New("ip not found")
	}
	return res
}

type PingResult struct {
	TargetAddr string
	RemoteAddr net.Addr
	Err        error
	Times      time.Duration
	Resp       *IcmpMessage
	Body       []byte
}

func ErrPingResult(targetAddr string, err error) *PingResult {
	return &PingResult{TargetAddr: targetAddr, Err: err}
}

func (p *PingResult) String() string {
	if p.Err == nil {
		return fmt.Sprintf("%d bytes from %s: times=%v", p.Resp.Len(), p.RemoteAddr.String(), p.Times)
	}
	return fmt.Sprintf("Error: ping %s, err: %s", p.TargetAddr, p.Err.Error())
}

func Pinger(address string, timeout time.Duration) *PingResult {
	c, err := net.Dial("ip4:icmp", address)
	if err != nil {
		return ErrPingResult(address, err)
	}
	err = c.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return ErrPingResult(address, err)
	}
	defer func() {
		_ = c.Close()
	}()
	typ := ICMPv4EchoRequest
	xid, xseq := os.Getpid()&0xffff, 1
	req := &IcmpMessage{
		Type: typ, Code: 0,
		Body: &ICMPEcho{
			ID: xid, Seq: xseq,
			Data: make([]byte, 56),
		},
	}
	wb, err := req.Marshal()
	if err != nil {
		return ErrPingResult(address, err)
	}
	start := time.Now()
	if _, err = c.Write(wb); err != nil {
		return ErrPingResult(address, err)
	}
	m := &IcmpMessage{}
	rb := make([]byte, 20+len(wb))
	var ttl time.Duration
	for {
		l, err := c.Read(rb)
		ttl = time.Since(start)
		rb = rb[:l]
		if err != nil {
			return ErrPingResult(address, err)
		}
		rb = ipv4Payload(rb)
		if err = m.Unmarshal(rb); err != nil {
			return ErrPingResult(address, err)
		}
		switch m.Type {
		case ICMPv4EchoRequest, ICMPv6EchoRequest:
			fmt.Println(ICMPv4EchoRequest)
			continue
		}
		break
	}
	return &PingResult{
		Err:        nil,
		Body:       rb,
		Resp:       m,
		Times:      ttl,
		TargetAddr: address,
		RemoteAddr: c.RemoteAddr(),
	}
}
func ipv4Payload(b []byte) []byte {
	if len(b) < 20 {
		return b
	}
	hdrlen := int(b[0]&0x0f) << 2
	return b[hdrlen:]
}
