package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/bigwhite/tcp-server-demo3/frame"
	"github.com/bigwhite/tcp-server-demo3/metrics"
	"github.com/bigwhite/tcp-server-demo3/packet"
)

func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(framePayload)
	if err != nil {
		fmt.Println("handleConn: packet decode error:", err)
		return
	}

	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		//fmt.Printf("recv submit: id = %s, payload=%s\n", submit.ID, string(submit.Payload))
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handleConn: packet encode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}

func handleConn(c net.Conn) {
	metrics.ClientConnected.Inc()
	defer func() {
		metrics.ClientConnected.Dec()
		c.Close()
	}()
	frameCodec := frame.NewMyFrameCodec()
	/*
		我们新增了一个读缓存变量（rbuf）和一个写缓存变量（wbuf），我们用这两个变量替换掉传给frameCodec.Decode和frameCodec.Encode的net.Conn参数。
	*/
	rbuf := bufio.NewReader(c)
	wbuf := bufio.NewWriter(c)

	defer wbuf.Flush()
	for {
		// read from the connection

		// decode the frame to get the payload
		// the payload is undecoded packet
		framePayload, err := frameCodec.Decode(rbuf)
		if err != nil {
			fmt.Println("handleConn: frame decode error:", err)
			return
		}
		metrics.ReqRecvTotal.Add(1)

		// do something with the packet
		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Println("handleConn: handle packet error:", err)
			return
		}

		// write ack frame to the connection
		err = frameCodec.Encode(wbuf, ackFramePayload)
		if err != nil {
			fmt.Println("handleConn: frame encode error:", err)
			return
		}
		metrics.RspSendTotal.Add(1)
	}
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	fmt.Println("server start ok(on *:8888)")

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go handleConn(c)
	}
}
