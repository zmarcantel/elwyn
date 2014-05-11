package chat

import (
	"net/http"
	"strconv"

	"code.google.com/p/go.net/websocket"

	"github.com/zmarcantel/elwyn/logging"
)

const (
	DEFAULT_ENTRANCE_BUFFER  int = 10
	DEFAULT_MESSAGE_BUFFER   int = 50
	DEFAULT_HEARTBEAT_BUFFER int = 100
)

var server *Server
var logger *logging.Router

var ErrClosedNetwork = "use of closed network connection"

type Server struct {
	Clients map[string]*Client
	Users   map[string]*Client
	Join    chan *Client
	Leave   chan *Client
	Receive chan *Message
	Error   chan error
}

func Initialize(lock chan error, log *logging.Router, port int) {
	logger = log

	server = &Server{
		Clients: make(map[string]*Client),
		Users:   make(map[string]*Client),
		Join:    make(chan *Client, DEFAULT_ENTRANCE_BUFFER),
		Leave:   make(chan *Client, DEFAULT_ENTRANCE_BUFFER),
		Receive: make(chan *Message, DEFAULT_MESSAGE_BUFFER),
		Error:   lock,
	}

	http.Handle("/chat", websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				server.Error <- err
			}
		}()

		logger.Println("Got chat connection")

		var client = NewClient(ws, server.Receive)
		server.Join <- client
		client.Listen(lock, server.Leave)
	}))

	logger.Banner("Starting Chat Server")
	server.Listen(lock, port)
}

func (self *Server) Listen(lock chan error, port int) {
	go func() {
		for {
			select {
			case joining := <-self.Join:
				self.Clients[joining.ID] = joining
				logger.Printf("New client: %s\n", joining.ID)
				joining.PingPong(lock, self)
				break

			case leaving := <-self.Leave:
				delete(self.Clients, leaving.ID)
				delete(self.Users, leaving.Username)
				logger.Printf("Signing off: %s\n", leaving.ID)
				self.AnnounceLeaveToRoom(leaving)
				break

			case data := <-self.Receive:
				if data.Action == "join" {
					var client = self.Clients[data.AuthorID]
					var res, created = checkNotExisting(self.Users, data.Author, client)
					if created {
						self.AnnounceJoinToRoom(data, client)
					}
					client.Write(lock, res)
					break
				}

				logger.Println("Broadcasting message")
				self.Broadcast(data)
				break

			case err := <-self.Error:
				lock <- err
				break
			}
		}
	}()

	lock <- http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (self *Server) Broadcast(message *Message) {
	for _, client := range self.Clients {
		message.Mine = (message.AuthorID == client.ID)
		client.Write(self.Error, message)
	}
}

func (self *Server) AnnounceJoinToRoom(data *Message, client *Client) {
	self.Broadcast(&Message{
		Action: "message",
		Author: "server",
		Body:   "<small><i>" + data.Author + "</i></small> has joined the room",
	})
}

func (self *Server) AnnounceLeaveToRoom(leaving *Client) {
	self.Broadcast(&Message{
		Action: "message",
		Author: "server",
		Body:   "<small><i>" + leaving.Username + "</i></small> has left the room",
	})
}

func checkNotExisting(users map[string]*Client, name string, client *Client) (*Message, bool) {
	if _, exists := users[name]; exists {
		return &Message{
			Author: "server",
			Action: "ACK",
			Body:   "exists",
		}, false
	}

	users[name] = client
	client.Username = name
	return &Message{
		Author: "server",
		Action: "ACK",
		Body:   "accepted",
	}, true
}
