package chat

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"image/png"
	"io"
	"time"

	sigil "github.com/cupcake/sigil/gen"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
)

const (
	DEFAULT_CLIENT_BUFFER int = 10
)

var identi = sigil.Sigil{
	Rows: 5,
	Foreground: []color.NRGBA{
		rgb(45, 79, 255),
		rgb(254, 180, 44),
		rgb(226, 121, 234),
		rgb(30, 179, 253),
		rgb(232, 77, 65),
		rgb(49, 203, 115),
		rgb(141, 69, 170),
	},
	Background: rgb(255, 255, 255),
}

func rgb(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 255} }

type Client struct {
	ID        string
	Username  string
	Socket    *websocket.Conn
	Receive   chan *Message
	Hangup    chan bool
	Heartbeat chan *Message
	Icon      string // base64 encoded
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
	var msg Message
	var err = websocket.JSON.Receive(self.Socket, &msg)

	if err != nil && err.Error() != ErrClosedNetwork {
		if err == io.EOF {
			self.Hangup <- true
			return
		}

		lock <- err
		return
	}
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

func (self *Client) Write(lock chan error, msg *Message) {
	var err = websocket.JSON.Send(self.Socket, msg)

	if err != nil && err.Error() != ErrClosedNetwork {
		logger.Errorln(err)
		lock <- err
		return
	}
}

func (self *Client) Ping(lock chan error) {
	self.Write(lock, &Message{
		Action: "heartbeat",
		Author: "server",
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
	var image = identi.Make(100, true, []byte(self.Username))
	var space = make([]byte, 0)
	var buf = bytes.NewBuffer(space)
	var result = base64.NewEncoder(base64.StdEncoding, buf)

	buf.Write([]byte("data:image/png;base64,"))
	png.Encode(result, image)
	self.Icon = buf.String()
}

func NewClient(sock *websocket.Conn, pipe chan *Message) *Client {
	if sock == nil {
		logger.Errorln("Received nil socket")
		return nil
	}

	return &Client{
		ID:        uuid.New(),
		Socket:    sock,
		Receive:   pipe,
		Hangup:    make(chan bool, 0), // blocks until read
		Heartbeat: make(chan *Message, 0),
	}
}
