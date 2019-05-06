package noise

import (
	"io"
	"sync"
	"time"
)

type Conn interface {
	io.Closer

	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}

type Dialer func(n *Node, address string) (*Peer, error)

type evt struct {
	opcode byte
	nonce  uint32
	msg    []byte

	oneway bool
	done   chan error
}

var evtPool sync.Pool

func acquireEvt() *evt {
	v := evtPool.Get()
	if v == nil {
		v = new(evt)
	}
	e := v.(*evt)
	if len(e.done) != 0 {
		panic("BUG: evt.done must be empty")
	}
	return e
}

func releaseEvt(e *evt) {
	if len(e.done) != 0 {
		panic("BUG: evt.done must be empty")
	}

	if e.done != nil {
		close(e.done)
		e.done = nil
	}
	evtPool.Put(e)
}
