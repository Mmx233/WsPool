package pool

import (
	"container/list"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

func New(upper *websocket.Upgrader) *Pool {
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
	el := a.List.Front()
	for el != nil {
		if !f(el.Value.(*Conn)) {
			break
		}
		el = el.Next()
	}
}

func (a *Pool) DoLoad(key any) (*Conn, bool) {
	var conn *Conn
	el := a.List.Front()
	for el != nil {
		if conn = el.Value.(*Conn); conn.Key == key {
			return conn, true
		}
		el = el.Next()
	}
	return nil, false
}

func (a *Pool) Load(key any) (*Conn, bool) {
	a.RLock()
	defer a.RUnlock()
	return a.DoLoad(key)
}

func (a *Pool) DoConnect(c *gin.Context, key any, resHeader http.Header) (*Conn, error) {
	ws, err := a.Upper.Upgrade(c.Writer, c.Request, resHeader)
	return &Conn{
		Conn: ws,
		Pool: a,
		Key:  key,
	}, err
}

func (a *Pool) NewConn(c *gin.Context, key any, resHeader http.Header) (*Conn, error) {
	conn, err := a.DoConnect(c, key, resHeader)
	if err != nil {
		return nil, err
	}
	a.Lock()
	defer a.Unlock()
	conn.Element = a.List.PushBack(conn)
	return conn, nil
}
