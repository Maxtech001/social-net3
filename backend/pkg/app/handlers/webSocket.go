package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = sendMessagefunc
}

func sendMessagefunc(event Event, c *Client) error {
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	var broadMessage NewMessageEvent
	var err error
	var conversation *models.Conversation

	broadMessage.Sent = time.Now().Format("2006-01-02 15:04:05")
	broadMessage.Message = chatevent.Message
	broadMessage.From = chatevent.From
	broadMessage.FromObj, err = dbfunctions.GetAccountDataById(broadMessage.From)
	if err != nil {
		return fmt.Errorf("failed to get sender data: %v", err)
	}
	broadMessage.To = chatevent.To
	broadMessage.GroupId = chatevent.GroupId
	broadMessage.IsRead = false

	if broadMessage.GroupId == 0 {
		conversation, err = dbfunctions.GetConversationBySenderReceiverIDs(broadMessage.From, broadMessage.To)
		if err != nil {
			return fmt.Errorf("failed to get conversation data: %v", err)
		}
	} else {
		conversation, err = dbfunctions.GetConversationByGroupId(broadMessage.GroupId)
		if err != nil {
			return fmt.Errorf("failed to get conversation data: %v", err)
		}
	}
	if conversation == nil {
		if broadMessage.GroupId == 0 {
			conversation, err = dbfunctions.InsertConversation(models.Conversation{User1ID: broadMessage.From, User2ID: broadMessage.To, GroupId: broadMessage.GroupId, CreatedAt: broadMessage.Sent})
			if err != nil {
				return fmt.Errorf("failed to insert conversation data: %v", err)
			}
		} else {
			conversation, err = dbfunctions.InsertConversation(models.Conversation{User1ID: 0, User2ID: 0, GroupId: broadMessage.GroupId, CreatedAt: broadMessage.Sent})
			if err != nil {
				return fmt.Errorf("failed to insert conversation data: %v", err)
			}
		}
	}

	message, err := dbfunctions.InsertMessage(models.Message{ConversationId: conversation.Id, SenderId: broadMessage.From, Message: broadMessage.Message, CreatedAt: broadMessage.Sent, IsRead: broadMessage.IsRead})
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	broadMessage.MessageId = message.Id

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	outgoingEvent := Event{
		Payload: data,
		Type:    EventNewMessage,
	}

	if broadMessage.GroupId == 0 {
		sendEvent(c, broadMessage.From, outgoingEvent)
		sendEvent(c, broadMessage.To, outgoingEvent)
	} else {
		users, err := dbfunctions.GetGroupMembersByGroupId(broadMessage.GroupId)
		if err != nil {
			fmt.Println("Getting group members failed: " + err.Error())
			return nil
		}
		for _, user := range users {
			sendEvent(c, user.Id, outgoingEvent)
		}
	}

	return nil
}

func sendEvent(c *Client, userId int, outgoingEvent Event) {
	if _, ok := c.manager.clients[userId]; ok {
		c.manager.clients[userId].egress <- outgoingEvent
	}
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request) {

	log.Println("New connection")

	user := AuthenticationService.GetUserFromCookie(r)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, m, user)

	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client.userId] = client
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if m.clients[client.userId] == client {
		client.connection.Close()
		delete(m.clients, client.userId)
	}
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:8080":
		fallthrough
	case "http://127.0.0.1:8080":
		return true
	default:
		return false
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	c.connection.SetReadLimit(512)

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.connection.SetPongHandler(c.pongHandler)
	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading messages: %v", err)
			}
			break
		}
		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handeling Message: ", err)
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
			log.Println("sent message")
		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return
			}
		}
	}
}

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[int]*Client

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	userId     int
}

func NewClient(conn *websocket.Conn, manager *Manager, user *models.User) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		userId:     user.Id,
	}
}

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
)

type SendMessageEvent struct {
	Message   string       `json:"message"`
	From      int          `json:"from"`
	FromObj   *models.User `json:"fromObj"`
	To        int          `json:"to"`
	GroupId   int          `json:"groupId"`
	MessageId int          `json:"messageId"`
	IsRead    bool         `json:"isRead"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent string `json:"sent"`
}
