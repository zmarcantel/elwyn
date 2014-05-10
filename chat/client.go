package chat

import (
	"io"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.net/websocket"
)

const (
	DEFAULT_CLIENT_BUFFER int = 10
)

type Client struct {
	ID       string
	Username string
	Socket   *websocket.Conn
	Receive  chan *Message
	Hangup   chan bool
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

	if err != nil {
		if err == io.EOF {
			self.Hangup <- true
			return
		}

		lock <- err
		return
	}

	logger.Printf("Received message from: %s\n", self.ID)
	msg.AuthorID = self.ID
	self.Receive <- &msg
}

func (self *Client) Write(lock chan error, msg *Message) {
	var err = websocket.JSON.Send(self.Socket, *msg)

	if err == nil {
		// TODO: better casing
	} else if err != nil {
		if err.Error() != ErrClosedNetwork {
			logger.Errorln("Attempted to use closed socket")
			logger.Warnln(err)
		}

		lock <- err
		return
	}
}

func NewClient(sock *websocket.Conn, pipe chan *Message) *Client {
	if sock == nil {
		logger.Errorln("Received nil socket")
		return nil
	}

	return &Client{
		ID:      uuid.New(),
		Socket:  sock,
		Receive: pipe,
		Hangup:  make(chan bool, 0), // blocks until read
	}
}
