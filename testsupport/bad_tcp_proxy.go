package testsupport

import (
	"errors"
	"io"
	"log"
	"net"
	"time"
)

func NewBadTCPProxy(destAddr string) (*BadTCPProxy, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	proxy := &BadTCPProxy{
		listener: listener,
		destAddr: destAddr,
	}

	go proxy.run()

	return proxy, nil
}

type BadTCPProxy struct {
	listener net.Listener
	destAddr string
}

func (bp *BadTCPProxy) run() error {
	defer bp.Close()

	for {
		clientConn, err := bp.listener.Accept()
		if err != nil {
			return err
		}

		go bp.handleClientConn(clientConn)
	}
}

func (bp *BadTCPProxy) handleClientConn(clientConn net.Conn) {
	connToDest, err := (&net.Dialer{}).Dial("tcp", bp.destAddr)
	if err != nil {
		log.Printf("%s", err)
		return
	}

	// dest receives request
	go io.Copy(connToDest, clientConn)
	// go bp.badCopy(connToDest, clientConn)

	// client receives response
	// go io.Copy(clientConn, connToDest)
	go bp.badCopy(clientConn, connToDest)
}

func (bp *BadTCPProxy) badCopy(dst io.Writer, src io.Reader) {
	sbw := &SingleByteWriter{
		Destination: dst,
		ByteInspector: func(b byte) error {
			if b == byte('O') {
				select {
				case <-time.After(time.Hour):
					return errors.New("potato timeout")
				}
			}
			return nil
		},
	}
	io.Copy(sbw, src)
}

func (bp *BadTCPProxy) Close() {
	bp.listener.Close()
}

func (bp *BadTCPProxy) ListenAddress() string {
	return bp.listener.Addr().String()
}
