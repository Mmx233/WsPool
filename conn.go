package pool

import (
	"container/list"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Conn struct {
	*websocket.Conn
	sync.Mutex
	Pool    *Pool
	Key     any
	onClose func(conn *Conn)
	Element *list.Element
}

func (a *Conn) WriteJSON(v interface{}) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WriteJSON(v)
}

func (a *Conn) WriteMessage(messageType int, data []byte) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WriteMessage(messageType, data)
}

func (a *Conn) WriteControl(messageType int, data []byte, deadline time.Time) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WriteControl(messageType, data, deadline)
}

func (a *Conn) WritePreparedMessage(pm *websocket.PreparedMessage) error {
	a.Lock()
	defer a.Unlock()
	return a.Conn.WritePreparedMessage(pm)
}

func (a *Conn) Clear() error {
	defer func() {
		a.Lock()
		if a.onClose != nil {
			a.onClose(a)
		}
		a.Unlock()
	}()
	a.Lock()
	ok := a.Element.Value == nil
	a.Unlock()
	if !ok {
		a.Pool.Lock()
		a.Pool.List.Remove(a.Element)
		a.Pool.Unlock()
		a.Lock()
		a.Element.Value = nil
		a.Unlock()
		return a.Close()
	}
	return nil
}

func (a *Conn) OnClose(e func(conn *Conn)) {
	a.Lock()
	defer a.Unlock()
	a.onClose = e
}

func (a *Conn) InPool() bool {
	a.Lock()
	defer a.Unlock()
	return a.Element != nil
}
