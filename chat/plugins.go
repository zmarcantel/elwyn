package chat

import (
    "time"

    "github.com/zmarcantel/elwyn/chat/common"
)

type Database interface {
    Write(msg *common.Message)      error
    LoadSince(time.Time)            []*common.Message
    LoadLast(count int64)           []*common.Message
    Close()                         error
}
