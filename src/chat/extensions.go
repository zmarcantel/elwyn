package chat

import (
	"time"

	"./common"
)

type Database interface {
	Write(msg *common.Packet) error
	LoadSince(time.Time) []*common.Packet
	LoadLast(count int64) []*common.Packet
	Close() error
}
