package netutils

import "errors"

const (
	ICMPv4EchoRequest = 8
	ICMPv4EchoReply   = 0
	ICMPv6EchoRequest = 128
	ICMPv6EchoReply   = 129
)

type IcmpMessage struct {
	Type     int       // type
	Code     int       // code
	Checksum int       // checksum
	Body     *ICMPEcho // body
}

func (m *IcmpMessage) Len() int {
	return 4 + m.Body.Len()
}

type Message interface {
	Len() int
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// Marshal returns the binary encoding of the ICMP echo request or
// reply message m.
func (m *IcmpMessage) Marshal() ([]byte, error) {
	b := []byte{byte(m.Type), byte(m.Code), 0, 0}
	if m.Body != nil && m.Body.Len() != 0 {
		mb, err := m.Body.Marshal()
		if err != nil {
			return nil, err
		}
		b = append(b, mb...)
	}
	switch m.Type {
	case ICMPv6EchoRequest, ICMPv6EchoReply:
		return b, nil
	}
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	b[2] ^= byte(^s & 0xff)
	b[3] ^= byte(^s >> 8)
	return b, nil
}

// Unmarshal parses b as an ICMP message.
func (m *IcmpMessage) Unmarshal(b []byte) error {
	msgLen := len(b)
	if msgLen < 4 {
		return errors.New("message too short")
	}
	m.Type = int(b[0])
	m.Code = int(b[1])
	m.Checksum = int(b[2])<<8 | int(b[3])
	if msgLen > 4 {
		var err error
		switch m.Type {
		case ICMPv4EchoRequest, ICMPv4EchoReply, ICMPv6EchoRequest, ICMPv6EchoReply:
			m.Body = new(ICMPEcho)
			err = m.Body.Unmarshal(b[4:])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ICMPEcho represenets an ICMP echo request or reply message body.
type ICMPEcho struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

func (p *ICMPEcho) Len() int {
	if p == nil {
		return 0
	}
	return 4 + len(p.Data)
}

// Marshal returns the binary enconding of the ICMP echo request or
// reply message body p.
func (p *ICMPEcho) Marshal() ([]byte, error) {
	b := make([]byte, 4+len(p.Data))
	b[0], b[1] = byte(p.ID>>8), byte(p.ID&0xff)
	b[2], b[3] = byte(p.Seq>>8), byte(p.Seq&0xff)
	copy(b[4:], p.Data)
	return b, nil
}

// Unmarshal parses b as an ICMP echo request or reply message body.
func (p *ICMPEcho) Unmarshal(b []byte) error {
	bodylen := len(b)
	p.ID = int(b[0])<<8 | int(b[1])
	p.Seq = int(b[2])<<8 | int(b[3])
	if bodylen > 4 {
		p.Data = make([]byte, bodylen-4)
		copy(p.Data, b[4:])
	}
	return nil
}
