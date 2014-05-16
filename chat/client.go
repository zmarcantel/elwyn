package chat

import (
	"bytes"
	"encoding/base64"
	"io"
	"time"

    "github.com/zmarcantel/elwyn/chat/common"

    "github.com/dgryski/go-identicon"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
)

const (
	DEFAULT_CLIENT_BUFFER int = 10
)

var identiconSalt = []byte{
    0x00, 0x11, 0x22, 0x33,
    0x44, 0x55, 0x66, 0x77,
    0x88, 0x99, 0xAA, 0xBB,
    0xCC, 0xDD, 0xEE, 0xFF,
}

type Client struct {
	ID        string
	Username  string
	Socket    *websocket.Conn
	Receive   chan *common.Message
	Hangup    chan bool
	Heartbeat chan *common.Message
	Icon      string // string of base64 encoded image
}

func (self *Client) Listen(lock chan error, done chan *Client) {
	go func() {
		for {
			hangup := <-self.Hangup
			if hangup {
				done <- self
				return
			}
		}
	}()

	for {
		self.Read(lock)
	}
}

func (self *Client) Read(lock chan error) {
	var msg common.Message
	var err = websocket.JSON.Receive(self.Socket, &msg)

	if err != nil && err.Error() != ErrClosedNetwork {
		if err == io.EOF {
			self.Hangup <- true
			return
		}

		lock <- err
		return
	}

    msg.Time = time.Now()
	msg.AuthorID = self.ID
	msg.Icon = self.Icon

	switch msg.Action {
	case "heartbeat":
		logger.Printf("Received Heartbeat: %s\n", self.ID)
		self.Heartbeat <- &msg
		break

	default:
		logger.Printf("Received message from: %s\n", self.ID)
		self.Receive <- &msg
		break
	}
}

func (self *Client) Write(lock chan error, msg *common.Message) {
	var err = websocket.JSON.Send(self.Socket, msg)

	if err != nil && err.Error() != ErrClosedNetwork {
		logger.Errorln(err)
		lock <- err
		return
	}
}

func (self *Client) Ping(lock chan error) {
	self.Write(lock, &common.Message{
		Action: "heartbeat",
		Sender: "server",
		Body:   "ping",
	})
}

func (self *Client) PingPong(lock chan error, serv *Server) {
	go func(server *Server) {
		for _ = range time.Tick(5 * time.Second) {
			self.Ping(lock)

			msg := <-self.Heartbeat

			if msg.Body != "pong" {
				logger.Errorf("Invalid Heartbeat From: %s\n\tReceived: %s\n\tExpected: pong\n", self.ID, msg.Body)
				server.Leave <- self
			}
		}
	}(serv)
}

func (self *Client) GenerateIcon() {
	var buf = bytes.NewBuffer(make([]byte, 0))
	var encoded = base64.NewEncoder(base64.StdEncoding, buf)

	buf.Write([]byte("data:image/png;base64,"))
    icon := identicon.New5x5(identiconSalt)
    pngdata := icon.Render([]byte(self.Username))
    encoded.Write(pngdata)
    self.Icon = buf.String()
}

func NewClient(sock *websocket.Conn, pipe chan *common.Message) *Client {
	if sock == nil {
		logger.Errorln("Received nil socket")
		return nil
	}

	return &Client{
		ID:        uuid.New(),
		Socket:    sock,
		Receive:   pipe,
		Hangup:    make(chan bool, 0), // blocks until read
		Heartbeat: make(chan *common.Message, 0), // TODO: test conditions of blocking and missing heartbeat
	}
}
