package chat

import (
	"net/http"
	"strconv"

	"code.google.com/p/go.net/websocket"

	"github.com/zmarcantel/elwyn/logging"
	"github.com/zmarcantel/elwyn/chat/common"
	"github.com/zmarcantel/elwyn/chat/databases/file"
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
	Receive chan *common.Message
	Error   chan error
    backing Database
}

func Initialize(lock chan error, log *logging.Router, port int)  Database {
	logger = log
    var fileBacking, err = file.Open("./messages.store")
    if err != nil {
        lock <- err
        return nil
    }

	server = &Server{
		Clients: make(map[string]*Client),
		Users:   make(map[string]*Client),
		Join:    make(chan *Client, DEFAULT_ENTRANCE_BUFFER),
		Leave:   make(chan *Client, DEFAULT_ENTRANCE_BUFFER),
		Receive: make(chan *common.Message, DEFAULT_MESSAGE_BUFFER),
		Error:   lock,
        backing: fileBacking,
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
	go server.Listen(lock, port)
    return fileBacking
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
                if (leaving == nil || leaving.Username == "") {
                    continue
                }
				delete(self.Clients, leaving.ID)
				delete(self.Users, leaving.Username)
				logger.Printf("Signing off: %s\n", leaving.ID)
				self.AnnounceLeaveToRoom(leaving)
				break

			case data := <-self.Receive:
				if data.Action == "join" {
                    self.HandleJoin(data)
					continue
				}
                self.backing.Write(data)
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

func (self *Server) Broadcast(message *common.Message) {
	for _, client := range self.Clients {
		message.Mine = (message.AuthorID == client.ID)
		client.Write(self.Error, message)
	}
}

func (self *Server) AnnounceJoinToRoom(data *common.Message, client *Client) {
	self.Broadcast(&common.Message{
		Action: "message",
		Sender: "server",
		Body:   "<small><i>" + data.Sender + "</i></small> has joined the room",
	})
}

func (self *Server) HandleJoin(data *common.Message) {
    var client = self.Clients[data.AuthorID]
    var res, created = checkNotExisting(self.Users, data.Sender, client)
    if created {
        self.AnnounceJoinToRoom(data, client)
        client.GenerateIcon()
    }
    client.Write(self.Error, res)
}

func (self *Server) AnnounceLeaveToRoom(leaving *Client) {
	self.Broadcast(&common.Message{
		Action: "message",
		Sender: "server",
		Body:   "<small><i>" + leaving.Username + "</i></small> has left the room",
	})
}

func checkNotExisting(users map[string]*Client, name string, client *Client) (*common.Message, bool) {
	if _, exists := users[name]; exists {
		return &common.Message{
			Sender: "server",
			Action: "ACK",
			Body:   "exists",
		}, false
	}

	users[name] = client
	client.Username = name
	return &common.Message{
		Sender: "server",
		Action: "ACK",
		Body:   "accepted",
	}, true
}

