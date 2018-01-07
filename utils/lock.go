package utils

type empty struct {}

type Lock struct {
	semaphore chan empty
}



func NewLock()*Lock{
	return &Lock{
		semaphore: make(chan empty, 1),
	}
}

func (lock *Lock) Acquire() {
	lock.semaphore <- empty{}
}

func (lock *Lock) Release(){
	<- lock.semaphore
}
