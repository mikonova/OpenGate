package async

import (
	"net"
	"sync"
	"time"
)

type SafeConn struct {
	c   net.Conn
	mut sync.RWMutex
}

func MakeSafeConn(conn net.Conn) *SafeConn {
	return &SafeConn{
		c: conn,
	}
}

func (sc *SafeConn) Read(src []byte) {
	sc.mut.RLock()
	sc.c.Read(src)
	defer sc.mut.RUnlock()
}

func (sc *SafeConn) Write(dst []byte) (int, error) {
	sc.mut.Lock()
	n, err := sc.c.Write(dst)
	defer sc.mut.Unlock()
	return n, err
}

// should not be called manually
func (sc *SafeConn) Close() {
	sc.c.Close()
}

func (sc *SafeConn) SetDeadline(t time.Duration) {
	sc.mut.Lock()
	sc.c.SetDeadline(time.Now().Add(t))
	sc.mut.Unlock()
}

func (sc *SafeConn) SetWriteDeadline(t time.Duration) {
	sc.mut.Lock()
	sc.c.SetWriteDeadline(time.Now().Add(t))
	sc.mut.Unlock()
}

func (sc *SafeConn) SetReadDeadline(t time.Duration) {
	sc.mut.Lock()
	sc.c.SetReadDeadline(time.Now().Add(t))
	sc.mut.Unlock()
}
