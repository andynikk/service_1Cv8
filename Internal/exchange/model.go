package exchange

import (
	"Service_1Cv8/internal/constants"
	"time"
)

type Queue struct {
	Name         string
	RoutingKey   string
	ExchangePlan string
	Number       int64
	Messages     []MessageBox
}

type MessageBox struct {
	Sender   string
	UID      string
	Priority int
	Type     string
	Object   string
	Massage  string
}

type AnswerMessage struct {
	Name         string
	RoutingKey   string
	ExchangePlan string
	Number       int64
	UID          string
	Type         string
	Object       string
	Massage      string
}

type MessageBolt struct {
	Bucket     string
	RoutingKey string
	Priority   string
	UID        string
	Type       string
	Object     string
	Massage    string
}

type MessagePudge struct {
	Pudge      string
	RoutingKey string
	Key        string
	Priority   string
	UID        string
	Type       string
	Object     string
	Message    string
}

type MessageKey struct {
	Key string
}

type MsgWithSorter struct {
	Sorter  string
	Message []byte
}

type ActionsWithQueue struct {
	Name     string
	Create   bool
	Del      bool
	Clearing bool
	TypeQ    string
}

type ExchangeQueueInfo struct {
	Name        string
	UploadAt    time.Time
	DownloadAt  time.Time
	TypeStorage string
	Size        string
	Messages    string
}

type PriorityMessages map[int][]MsgWithSorter

type ExchangeQueue struct {
	Name               string
	UploadAt           time.Time
	DownloadAt         time.Time
	TypeStorage        constants.TypeStorage
	Size               int64
	TimeTransferToDisk float64
	MaxSizeSoftMode    int
	MaxLenSoftMode     int
	PriorityMessages
}

type ExchangeStorage map[string]ExchangeQueue
