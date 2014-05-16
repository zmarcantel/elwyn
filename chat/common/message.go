package common

import (
    "time"
)

type Message struct {
    // message data
	Sender   string     `json:"sender"`
	Receiver string     `json:"receiver"`
	Body     string     `json:"body"`
    Time     time.Time  `json:"timestamp"`

    // metadata -- descriptors to funnel messages
	Mine     bool       `json:"mine"`
	Action   string     `json:"action"`

    // sender data
	Icon     string     `json:"icon"`
	AuthorID string
}

func (self *Message) String() string {
	return self.Sender + " says " + self.Body
}
