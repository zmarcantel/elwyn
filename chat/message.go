package chat

type Message struct {
	Author   string `json:"author"`
	AuthorID string
	Body     string `json:"body"`
	Mine     bool   `json:"mine"`
	Action   string `json:"action"`
}

func (self *Message) String() string {
	return self.Author + " says " + self.Body
}
