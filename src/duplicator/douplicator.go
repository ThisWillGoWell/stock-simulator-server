package duplicator

import (
	"fmt"
	"time"

	"encoding/json"
	"reflect"

	"github.com/ThisWillGoWell/stock-simulator-server/src/deepcopy"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
)

/**
Channel Duplicators are a core part of the design and are used to link together the pipeline of the program
They basically just a fan-in -> fan-out channel
inputs get fanned onto the transfer channel, then out to each listening.
This makes linking different parts of the program together super easy since you can tie outputs to inputs.

*/
type ChannelDuplicator struct {
	transfer  chan interface{}
	outputs   []chan interface{}
	inputs    []<-chan interface{}
	debug     bool
	debugName string
	copy      bool
	sink      bool
	lock      *lock.Lock
	close     chan interface{}
	closed    chan interface{}
}

/**
name is only used for debugging, that's honestly the worst part of this design choice
Tracing out a rouge message is impossible without this (cant really break on a line of code being called
multiple times for each event)

*/
func MakeDuplicator(name string) *ChannelDuplicator {
	chDoup := &ChannelDuplicator{
		lock:      lock.NewLock("channel-duplicator"),
		outputs:   make([]chan interface{}, 0),
		inputs:    make([]<-chan interface{}, 0),
		transfer:  make(chan interface{}, 100),
		debug:     false,
		copy:      false,
		debugName: name,
		sink:      false,
		close:     make(chan interface{}),
		closed:    make(chan interface{}),
	}
	chDoup.startDuplicator()
	return chDoup
}

/**
Copy mode is
*/
func (ch *ChannelDuplicator) EnableCopyMode() {
	ch.copy = true
}

func (ch *ChannelDuplicator) EnableSink() {
	ch.copy = true
}

func (ch *ChannelDuplicator) DiableSink() {
	ch.copy = true
}

/**
a lot of printing but uuids make is not to bad to trace though whats happening
*/
func (ch *ChannelDuplicator) EnableDebug() {
	ch.debug = true
}

func (ch *ChannelDuplicator) SetName(name string) {
	ch.debugName = name
}

/**
Return a new output channel that is being fan-out from the transfer
*/
func (ch *ChannelDuplicator) GetOutput() chan interface{} {
	ch.lock.Acquire("getOutput")
	defer ch.lock.Release()
	// make a channel with a 10 buffer size
	if ch.debug {
		fmt.Println("adding output on", ch.debugName)
	}
	newOutput := make(chan interface{}, 100)
	ch.outputs = append(ch.outputs, newOutput)
	return newOutput
}

/**
Return a new buffered output
essentially used where ever you can have different levels of processing at each end
(a socket would need a bufferd since you can add messages quicker than you can send them)

*/
func (ch *ChannelDuplicator) GetBufferedOutput(buffSize int64) chan interface{} {
	// make a channel with a 10 buffer size
	if ch.debug {
		fmt.Println("adding output on", ch.debugName)
	}
	newOutput := make(chan interface{}, buffSize)
	ch.outputs = append(ch.outputs, newOutput)
	return newOutput
}

/**
remove a channel
*/
func (ch *ChannelDuplicator) UnregisterOutput(remove chan interface{}) {
	ch.lock.Acquire("remove-output")
	defer ch.lock.Release()
	removeIndex := -1
	for i, channel := range ch.outputs {
		if channel == remove {
			removeIndex = i
		}
	}
	if removeIndex == -1 {
		fmt.Println("can't find chan in unregister output")
		return
	}
	//Remove channel by swapping the removed channel to the end and then just trimming the slice
	ch.outputs[len(ch.outputs)-1], ch.outputs[removeIndex] = ch.outputs[removeIndex], ch.outputs[len(ch.outputs)-1]
	ch.outputs = ch.outputs[:len(ch.outputs)-1]
}

/**
Register a input to be fanned onto the transfer
*/
func (ch *ChannelDuplicator) RegisterInput(inputChannel <-chan interface{}) {
	ch.inputs = append(ch.inputs, inputChannel)
	go func() {
		if ch.debug {
			fmt.Println("registering input on", ch.debugName)
		}
		for val := range inputChannel {
			if ch.debug {
				str, _ := json.Marshal(val)
				fmt.Println("adding to trasfer=", ch.debugName, "value=", string(str))
			}
			ch.Offer(val)
			if ch.debug {
				fmt.Println("done transfer on", ch.debugName)
			}

		}
		if ch.debug {
			fmt.Println("closeing input on", ch.debugName)
		}
	}()

}

/**
offer a single value directly onto the transfer
*/
func (ch *ChannelDuplicator) Offer(value interface{}) {
	if ch.debug {
		fmt.Println("offering to transfer", ch.debugName)
	}
	if ch.copy && (reflect.TypeOf(value).Kind() == reflect.Ptr || reflect.TypeOf(value).Kind() == reflect.Slice) {
		newVal := deepcopy.Copy(value)
		//pass that pointer down the transfer line
		if ch.debug {
			str, _ := json.Marshal(newVal)
			fmt.Println("offering copy to transfer=", ch.debugName, "value=", string(str))
		}
		ch.transfer <- newVal
	} else {
		if ch.debug {
			str, _ := json.Marshal(value)
			fmt.Println("offering to trasfer=", ch.debugName, "value=", string(str))
		}
		ch.transfer <- value
	}
}

/**
Start running the fan out from transfer to all outputs
*/
func (ch *ChannelDuplicator) startDuplicator() {
	go func() {
		run := true
		for run {
			select {
			case nextValue := <-ch.transfer:
				ch.lock.Acquire("startDuplicator")
				if ch.debug {
					fmt.Println("sending down outputs on", ch.debugName)
				}
				if !ch.sink {
					for i, channel := range ch.outputs {
						select {
						case channel <- nextValue:
							if ch.debug {
								str, _ := json.Marshal(nextValue)
								fmt.Println("sent to an output of", ch.debugName, "index", i, "vaule", string(str))
							}
							continue
						default:
							fmt.Println("missing messages on", ch.debugName, "index", i)
							continue
						}
					}
				}
				ch.lock.Release()
			case <-ch.close:
				run = false
			}
		}
		close(ch.transfer)
		for _, ch := range ch.outputs {
			close(ch)
		}
	}()

}
func (ch *ChannelDuplicator) StopDuplicator() {
	ch.close <- nil
}

/**
Unlink two duplicators from each other
*/
func UnlinkDouplicator(input, output *ChannelDuplicator) {
	for _, inputCh := range input.inputs {
		for _, outputCh := range output.outputs {
			if inputCh == outputCh {
				fmt.Println("unlinked")
				output.UnregisterOutput(outputCh)
				close(outputCh)
			}
		}
	}
}

/**
test function
will print:

received: hello
received: hello
received: hello
received: world
received: world
received: world

*/
func main() {
	input1 := make(chan interface{})
	input2 := make(chan interface{})

	chDoup := MakeDuplicator("'")
	chDoup.RegisterInput(input1)
	chDoup.RegisterInput(input2)
	for i := 0; i < 3; i++ {
		output := chDoup.GetOutput()
		go func() {
			for value := range output {
				fmt.Println("recieved: ", value)
			}
		}()
	}
	input1 <- "hello"
	input2 <- "world"
	time.Sleep(time.Second * 1)
}
