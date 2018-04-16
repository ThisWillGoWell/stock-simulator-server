package lock

import "fmt"

type empty struct{}

/**
Not sure why, but when I first looked for locks in golang, I was told
Just use a channel!, so i did. I have not found a reason to go back and use
the actual sync library since this actually works well and is just super simple
a buffered channel can act as a semaphore where acquiring is writing to the channel
and release is reading
*/
type Lock struct {
	name      string
	semaphore chan empty
	debug     bool
}

/**
Locks have a name so they can be easily debugged when there is a deadlock
and I forget to release it some where
*/
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

/**
Loc is also used for debugging
*/
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
