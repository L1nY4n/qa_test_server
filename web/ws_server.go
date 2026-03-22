package web

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var WsManager = clientManager{
	clientGroup: make(map[string]map[string]*wsClient),
	register:    make(chan *wsClient),
	unRegister:  make(chan *wsClient),
	broadcast:   make(chan *boradcastData, 10),
}

type clientManager struct {
	clientGroup map[string]map[string]*wsClient
	register    chan *wsClient
	unRegister  chan *wsClient
	broadcast   chan *boradcastData

	groupCount  int64
	clientCount int64
}

type boradcastData struct {
	GroupID string
	Data    []byte
}

type wsClient struct {
	ID     string
	Group  string
	Socket *websocket.Conn
	Send   chan []byte
}

type WsStats struct {
	Groups  int64 `json:"groups"`
	Clients int64 `json:"clients"`
}

func (c *wsClient) Read() {
	defer func() {
		WsManager.unRegister <- c
		_ = c.Socket.Close()
	}()

	for {
		if _, _, err := c.Socket.ReadMessage(); err != nil {
			break
		}
	}
}

func (c *wsClient) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.Socket.WriteMessage(websocket.BinaryMessage, message)
		}
	}
}

func (manager *clientManager) Start() {
	log.Println("Websocket manager started")
	for {
		select {
		case client := <-manager.register:
			log.Printf("Websocket client %s connect", client.ID)
			if manager.clientGroup[client.Group] == nil {
				manager.clientGroup[client.Group] = make(map[string]*wsClient)
				atomic.AddInt64(&manager.groupCount, 1)
			}
			manager.clientGroup[client.Group][client.ID] = client
			atomic.AddInt64(&manager.clientCount, 1)
			log.Printf("Register client %s to group %s success", client.ID, client.Group)

		case client := <-manager.unRegister:
			log.Printf("Unregister websocket client %s", client.ID)
			if _, ok := manager.clientGroup[client.Group]; ok {
				if _, ok := manager.clientGroup[client.Group][client.ID]; ok {
					close(client.Send)
					delete(manager.clientGroup[client.Group], client.ID)
					atomic.AddInt64(&manager.clientCount, -1)
					log.Printf("Unregister websocket client %s from group %s success", client.ID, client.Group)
					if len(manager.clientGroup[client.Group]) == 0 {
						delete(manager.clientGroup, client.Group)
						atomic.AddInt64(&manager.groupCount, -1)
					}
				}
			}

		case data := <-manager.broadcast:
			if groupMap, ok := manager.clientGroup[data.GroupID]; ok {
				for _, conn := range groupMap {
					// Skip slow clients to keep broadcaster responsive.
					select {
					case conn.Send <- data.Data:
					default:
					}
				}
			}
		}
	}
}

func (manager *clientManager) RegisterClient(ctx *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("websocket client connect %v error", ctx.Param("channel"))
		return
	}

	client := &wsClient{
		ID:     uuid.NewV4().String(),
		Group:  ctx.Param("channel"),
		Socket: conn,
		Send:   make(chan []byte, 1024),
	}
	manager.register <- client
	go client.Read()
	go client.Write()
}

func (manager *clientManager) Groupbroadcast(group string, message []byte) {
	manager.broadcast <- &boradcastData{
		GroupID: group,
		Data:    message,
	}
}

func (manager *clientManager) TryGroupbroadcast(group string, message []byte) bool {
	select {
	case manager.broadcast <- &boradcastData{
		GroupID: group,
		Data:    message,
	}:
		return true
	default:
		return false
	}
}

func (manager *clientManager) Stats() WsStats {
	return WsStats{
		Groups:  atomic.LoadInt64(&manager.groupCount),
		Clients: atomic.LoadInt64(&manager.clientCount),
	}
}

func WsRoute(r *gin.Engine) {
	streamRoute := r.Group("/realtime")
	{
		streamRoute.GET("/register/:channel", Register)
	}
}

func Register(c *gin.Context) {
	WsManager.RegisterClient(c)
}
