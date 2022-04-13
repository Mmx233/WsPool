package pool

import (
	"container/list"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

func NewPool(upper *websocket.Upgrader) *Pool {
	return &Pool{
		List:  list.New(),
		Upper: upper,
	}
}

type Pool struct {
	sync.RWMutex
	Upper *websocket.Upgrader
	List  *list.List
}

func (a *Pool) Len() int {
	a.RLock()
	defer a.RUnlock()
	return a.List.Len()
}

func (a *Pool) Range(f func(*Conn) bool) {
	a.RLock()
	defer a.RUnlock()
	e := a.List.Front()
	for e != nil {
		if !f(e.Value.(*Conn)) {
			break
		}
		e = e.Next()
	}
}

func (a *Pool) Load(key any) (*Conn, bool) {
	a.RLock()
	defer a.RUnlock()
	var conn *Conn
	e := a.List.Front()
	for e != nil {
		if conn = e.Value.(*Conn); conn.Key == key {
			return conn, true
		}
		e = e.Next()
	}
	return nil, false
}

func (a *Pool) NewConn(c *gin.Context, key any, resHeader http.Header) (*Conn, error) {
	ws, e := a.Upper.Upgrade(c.Writer, c.Request, resHeader)
	if e != nil {
		return nil, e
	}
	a.Lock()
	defer a.Unlock()
	old, ok := a.Load(key)
	if ok {
		old.Lock()
		_ = old.Conn.Close()
		a.List.Remove(old.Element)
		old.Element.Value = nil
		old.Unlock()
	}
	conn := &Conn{
		Conn: ws,
		Pool: a,
		Key:  key,
	}
	conn.Element = a.List.PushBack(conn)
	return conn, nil
}
