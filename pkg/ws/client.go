package ws

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go_kratos_template/pkg/log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Options func(*Client)

const (
	Default = iota
	Running
	Exit
)

type Client struct {
	URL        string
	Conn       *websocket.Conn
	SendMsg    chan []byte
	ReceiveMsg chan []byte
	CloseChan  chan struct{}
	Status     int
	ID         string
}

func NewClient(opts ...Options) *Client {
	client := &Client{
		SendMsg:    make(chan []byte),
		ReceiveMsg: make(chan []byte),
		CloseChan:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func WithURL(url string) Options {
	return func(client *Client) {
		client.URL = url
	}
}

func WithUUID() Options {
	return func(client *Client) {
		client.ID = uuid.New().String()
	}
}

func WithID(id string) Options {
	return func(client *Client) {
		client.ID = id
	}
}

func (c *Client) UpGrader(w http.ResponseWriter, r *http.Request) error {
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		//设置请求协议
		Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(w, r, nil)
	c.Conn = conn
	return err
}

func (c *Client) Read() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[WebSocket] recover[%v] %v ", c.ID, err)
		}
	}()

	defer func() {
		c.Status = Exit
		c.Conn.Close()
		close(c.ReceiveMsg)
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Infof("[WebSocket] error: %v", err)
			}
			break
		}
		if len(message) > 0 {
			c.ReceiveMsg <- message
		}
	}
}

func (c *Client) Write() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[WebSocket] recover[%v] %v ", c.ID, err)
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.Status = Exit

		ticker.Stop()
		c.Conn.Close()
		close(c.SendMsg)
		close(c.CloseChan)
	}()

	for {
		select {
		case msg, ok := <-c.SendMsg:
			if !ok {
				//	chan closed
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Warnf("[WebSocket] CloseMessage Error: %v", err)
					return
				}
				log.Warnf("[WebSocket] SendMsg Chan Closed")
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Errorf("[WebSocket] WriteMessage Error: %v", err)
				return
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("[WebSocket] PingMessage Error: %v", err)
				return
			}
		case <-c.CloseChan:
			log.Errorf("[WebSocket] closeChan [%v]", c.ID)
			return
		}
	}
}

func (c *Client) Run() error {
	if c.Conn == nil {
		return fmt.Errorf("升级WebSocket协议失败")
	}
	c.Status = Running
	go c.Read()
	go c.Write()
	return nil
}
