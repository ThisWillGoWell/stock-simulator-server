package lock

type empty struct{}

type QueuedLock struct {
	name string
	lock *Lock
}

type Lock struct {
	name      string
	semaphore chan empty
}

func NewLock(name string) *Lock {
	return &Lock{
		name:      name,
		semaphore: make(chan empty, 1),
	}
}

func (lock *Lock) Acquire(loc string) {
	// fmt.Println("acquire: ", lock.name, "form:", loc)
	lock.semaphore <- empty{}
	// fmt.Println("acquired: ", lock.name, "from", loc)
}

func (lock *Lock) Release() {
	<-lock.semaphore
}
