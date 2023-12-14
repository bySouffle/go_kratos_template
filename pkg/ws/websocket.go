package ws

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var (
	once          sync.Once
	clientManager *ClientManager
)

type WebSocketHandler func(*websocket.Conn)

type ClientManager struct {
	group        map[string][]*Client //	URL分组
	registerChan chan *Client
	logOutChan   chan *Client
	log          *log.Helper
	ClientCnt    uint32
	closeChan    chan struct{}
	Status       int
}

func NewClientManager(logger log.Logger) *ClientManager {
	return &ClientManager{
		group:        make(map[string][]*Client),
		registerChan: make(chan *Client, 1),
		logOutChan:   make(chan *Client, 1),
		log:          log.NewHelper(logger),
		closeChan:    make(chan struct{}),
	}
}

func NewNewClientManagerWithRun(logger log.Logger) *ClientManager {
	once.Do(func() {
		clientManager = &ClientManager{
			group:        make(map[string][]*Client),
			registerChan: make(chan *Client, 1),
			logOutChan:   make(chan *Client, 1),
			log:          log.NewHelper(logger),
			closeChan:    make(chan struct{}),
		}
		clientManager.Start()
	})
	return clientManager
}

func (m *ClientManager) Register(client *Client) error {
	if m.Status == Exit {
		return fmt.Errorf("ClientManager已退出")
	}
	select {
	case m.registerChan <- client:
	default:
		return fmt.Errorf("register通道关闭")
	}
	return nil

}

func (m *ClientManager) register(client *Client) {
	m.group[client.URL] = append(m.group[client.URL], client)
	m.ClientCnt++
	m.log.Infof("[WebSocket] [%v]客户端注册:%v", client.URL, client.ID)
}

func (m *ClientManager) LogOut(client *Client) error {
	if m.Status == Exit {
		return fmt.Errorf("ClientManager已退出")
	}
	select {
	case m.logOutChan <- client:
		return nil
	default:
		return fmt.Errorf("logOut通道关闭")
	}
}

func (m *ClientManager) logOut(client *Client) {
	clients, found := m.group[client.URL]
	if !found {
		return
	}
	for i, c := range clients {
		if c.ID == client.ID {
			c.CloseChan <- struct{}{} //	关闭发送客户端
			m.group[client.URL] = append(clients[:i], clients[i+1:]...)
			m.ClientCnt--
			break
		}
	}
	m.log.Infof("[WebSocket] [%v]客户端注销:%v", client.URL, client.ID)
}

func (m *ClientManager) Start() {
	go func() {
		defer func() {
			m.Status = Exit
			if err := recover(); err != nil {
				log.Errorf("[WebSocket] recover[%v]", err)
			}
		}()
		defer close(m.registerChan)
		defer close(m.logOutChan)
		defer close(m.closeChan)
		m.Status = Running
		for {
			select {
			case client := <-m.registerChan:
				m.register(client)
			case client := <-m.logOutChan:
				m.logOut(client)
			case <-time.After(time.Second * 5):
				m.deleteClientWithExit()
			case <-m.closeChan:
				m.releaseGroup()
				m.log.Infof("[WebSocket] server stop")
				return
			}
		}
	}()
}
func (m *ClientManager) Stop() error {
	select {
	case m.closeChan <- struct{}{}:
	default:
		return fmt.Errorf("[WebSocket] closeChan closed")
	}
	return nil
}

func (m *ClientManager) releaseGroup() {
	var wg sync.WaitGroup
	for key, clients := range m.group {
		for _, client := range clients {
			wg.Add(1)
			go func(c *Client) {
				defer wg.Done()
				select {
				case c.CloseChan <- struct{}{}:
				default:
					log.Warnf("[WebSocket] %v CloseChan closed", c.ID)
				}
			}(client)
		}
		delete(m.group, key)
	}
	wg.Wait()
}

func (m *ClientManager) deleteClientWithExit() {
	for groupName, clients := range m.group {
		// 遍历组中的每个 Client 对象
		var RunningClient []*Client
		for i := 0; i < len(clients); i++ {
			// 如果状态为Exit，则删除该 Client 对象
			if clients[i].Status != Exit {
				RunningClient = append(RunningClient, clients[i])
			}
		}
		m.group[groupName] = RunningClient
	}
}

func (m *ClientManager) SendMsgToGroup(url string, msg []byte) {
	clients, found := m.group[url]
	if !found {
		return
	}
	wg := sync.WaitGroup{}
	for _, client := range clients {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			select {
			case c.SendMsg <- msg:
			default:
				log.Warnf("[WebSocket] %v SendMsg closed", c.ID)
			}
		}(client)
	}
	wg.Wait()
}

func (m *ClientManager) GetClientMsgChan(url string, id string) chan []byte {
	clients, found := m.group[url]
	if !found {
		return nil
	}
	for _, c := range clients {
		if c.ID == id {
			return c.ReceiveMsg
		}
	}
	return nil
}
