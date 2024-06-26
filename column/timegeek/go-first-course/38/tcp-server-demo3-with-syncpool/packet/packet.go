package packet

import (
	"bytes"
	"fmt"
	"sync"
)

// Packet协议定义

/*
### packet header
1 byte: commandID

### submit packet

8字节 ID 字符串
任意字节 payload

### submit ack packet

8字节 ID 字符串
1字节 result
*/

const (
	CommandConn   = iota + 0x01 // 0x01
	CommandSubmit               // 0x02
)

const (
	CommandConnAck   = iota + 0x80 // 0x81
	CommandSubmitAck               //0x82
)

type Packet interface {
	Decode([]byte) error     // []byte -> struct
	Encode() ([]byte, error) //  struct -> []byte
}

type Submit struct {
	ID      string
	Payload []byte
}

func (s *Submit) Decode(pktBody []byte) error {
	s.ID = string(pktBody[:8])
	s.Payload = pktBody[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (s *SubmitAck) Decode(pktBody []byte) error {
	s.ID = string(pktBody[0:8])
	s.Result = uint8(pktBody[8])
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), []byte{s.Result}}, nil), nil
}

/*
*
一提到重用内存对象，我们就会想到了sync.Pool。简单来说，sync.Pool就是官方实现的一个可复用的内存对象池，使用sync.Pool，我们可以减少堆对象分配的频度，进而降低给GC带去的压力。

我们继续在tcp-server-demo3的基础上，使用sync.Pool进行堆内存对象分配的优化，新版的代码放在了tcp-server-demo3-with-syncpool中。

新版代码相对于tcp-server-demo3有两处改动，第一处是在packet.go中，我们创建了一个SubmitPool变量，它的类型为sync.Pool，这就是我们的内存对象池，池中的对象都是Submit。这样我们在packet.Decode中收到Submit类型请求时，也不需要新分配一个Submit对象，而是直接从SubmitPool代表的Pool池中取出一个复用。
*/
var SubmitPool = sync.Pool{
	New: func() interface{} {
		return &Submit{}
	},
}

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	pktBody := packet[1:]

	switch commandID {
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandSubmit:
		//s := Submit{}
		s := SubmitPool.Get().(*Submit)
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return s, nil
	case CommandSubmitAck:
		s := SubmitAck{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

func Encode(p Packet) ([]byte, error) {
	var commandID uint8
	var pktBody []byte
	var err error

	switch t := p.(type) {
	case *Submit:
		commandID = CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{[]byte{commandID}, pktBody}, nil), nil
}
