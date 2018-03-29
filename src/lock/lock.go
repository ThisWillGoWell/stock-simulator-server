package lock

import "fmt"

type empty struct{}

type QueuedLock struct {
	name string
	lock *Lock
}

type Lock struct {
	name      string
	semaphore chan empty
	debug     bool
}

func NewLock(name string) *Lock {
	return &Lock{
		name:      name,
		semaphore: make(chan empty, 1),
		debug:     false,
	}
}

func (lock *Lock) EnableDebug() {
	lock.debug = true
}

func (lock *Lock) Acquire(loc string) {
	if lock.debug {
		fmt.Printf("1 get: %s, -> %s\n", lock.name, loc)
	}
	lock.semaphore <- empty{}
	if lock.debug {
		fmt.Printf("2 get: %s, -> %s\n", lock.name, loc)
	}
}

func (lock *Lock) Release() {
	if lock.debug {
		fmt.Printf("release: %s\n", lock.name)
	}
	<-lock.semaphore
}
