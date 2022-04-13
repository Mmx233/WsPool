package pool

import (
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
	if a.onClose != nil {
		defer a.onClose(a)
	}
	a.Pool.Lock()
	ws, ok := a.Pool.conn[a.Key]
	if ok && ws == a {
		delete(a.Pool.conn, a.Key)
		a.Pool.Unlock()
		a.Lock()
		defer a.Unlock()
		return a.Close()
	} else {
		a.Pool.Unlock()
	}
	return nil
}

func (a *Conn) OnClose(e func(conn *Conn)) {
	a.onClose = e
}
