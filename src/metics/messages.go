package metics

import "sync/atomic"

type MessageCount struct {
	Tx uint32
	Rx uint32
}

var counter = MessageCount{}

func MessageSent() {
	atomic.AddUint32(&counter.Tx, 1)
}

func MessageRecieved() {
	atomic.AddUint32(&counter.Rx, 1)
}
