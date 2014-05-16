package file

import (
    "os"
    "fmt"
    "bufio"
    "time"

    "github.com/zmarcantel/elwyn/chat/common"
)

type Database struct {
    path            string
    file            *os.File
    input           *bufio.Reader
    output          *bufio.Writer
}

func Open(path string) (result *Database, err error) {
    var desc *os.File
    desc, err = os.OpenFile(path, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0644)
    if err != nil { return }

    result = &Database{
        path:       path,
        file:       desc,
        input:      bufio.NewReader(desc),
        output:     bufio.NewWriter(desc),
    }

    return
}

func (self *Database) Write(msg *common.Message) (err error) {
    var out = fmt.Sprintf("%s | %s | %s\n", msg.Sender, msg.Time, msg.Body)
    _, err = self.output.WriteString(out)
    return
}

func (self *Database) LoadLast(count int64) []*common.Message {
    return []*common.Message {
        {
            Sender: "server",
            Action: "message",
            Body: "We're sorry, but the file backing does not currently support back-tracking",
        },
        {
            Sender: "server",
            Action: "message",
            Body: "Want to help? Submit a PR @ <a href=\"https://github.com/zmarcantel/elwyn\">the project page</a>",
        },
    }
}

func (self *Database) LoadSince(since time.Time) []*common.Message {
    return []*common.Message {
        {
            Sender: "server",
            Action: "message",
            Body: "We're sorry, but the file backing does not currently support back-tracking",
        },
        {
            Sender: "server",
            Action: "message",
            Body: "Want to help? Submit a PR @ <a href=\"https://github.com/zmarcantel/elwyn\">the project page</a>",
        },
    }
}

func (self *Database) Close() (err error) {
    err = self.output.Flush()
    return
}

func (self *Database) GetPath() string {
    return self.path
}

func (self *Database) GetFile() *os.File {
    return self.file
}

