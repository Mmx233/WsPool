package pool

import (
	"container/list"
	"github.com/gorilla/websocket"
	"io"
	"sync"
	"time"
)

type Conn struct {
	*websocket.Conn
	sync.Mutex
	Pool    *Pool
	Key     any
	OnClose func(conn *Conn)
	Element *list.Element
}

func (a *Conn) WriteReader(messageType int, r io.Reader) error {
	a.Lock()
	defer a.Unlock()
	w, err := a.Conn.NextWriter(messageType)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
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

func (a *Conn) DoClear() bool {
	if a.Element.Value == nil {
		return false
	}
	if a.OnClose != nil {
		defer a.OnClose(a)
	}
	a.Pool.Lock()
	a.Pool.List.Remove(a.Element)
	a.Pool.Unlock()
	a.Element.Value = nil
	_ = a.Close()
	return true
}

func (a *Conn) Clear() bool {
	a.Lock()
	defer a.Unlock()
	return a.DoClear()
}
