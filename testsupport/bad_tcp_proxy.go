package testsupport

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

func NewBadTCPProxy(destAddr string) (*BadTCPProxy, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	proxy := &BadTCPProxy{
		context:    ctx,
		cancelFunc: cancelFunc,
		listener:   listener,
		destAddr:   destAddr,
	}

	go proxy.run()

	return proxy, nil
}

type BadTCPProxy struct {
	context    context.Context
	cancelFunc context.CancelFunc
	listener   net.Listener
	destAddr   string
}

func (bp *BadTCPProxy) run() error {
	defer bp.Close()

	for {
		if bp.context.Err() != nil {
			return bp.context.Err()
		}

		clientConn, err := bp.listener.Accept()
		if err != nil {
			return err
		}

		go bp.handleClientConn(clientConn)
	}
}

func (bp *BadTCPProxy) handleClientConn(clientConn net.Conn) {
	connToDest, err := (&net.Dialer{}).DialContext(bp.context, "tcp", bp.destAddr)
	if err != nil {
		log.Printf("%s", err)
		return
	}

	go bp.badCopy(clientConn, connToDest)
	go bp.badCopy(connToDest, clientConn)
}

func (bp *BadTCPProxy) badCopy(dst io.Writer, src io.Reader) {
	sbw := &SingleByteWriter{
		Destination: dst,
		ByteInspector: func(b byte) error {
			if b == byte('O') {
				select {
				case <-time.After(time.Hour):
					return errors.New("potato timeout")
				case <-bp.context.Done():
					return bp.context.Err()
				}
			}
			return nil
		},
	}
	io.Copy(sbw, src)
}

func (bp *BadTCPProxy) Close() {
	bp.listener.Close()
	bp.cancelFunc()
}

func (bp *BadTCPProxy) ListenAddress() string {
	return bp.listener.Addr().String()
}
