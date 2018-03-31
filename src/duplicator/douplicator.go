package duplicator

import (
	"fmt"
	"time"

	"encoding/json"
	"github.com/stock-simulator-server/src/lock"
	"reflect"
	"github.com/stock-simulator-server/src/deepcopy"
)

type ChannelDuplicator struct {
	transfer  chan interface{}
	outputs   []chan interface{}
	inputs    []<-chan interface{}
	debug     bool
	debugName string
	copy      bool
	lock      *lock.Lock
}

func MakeDuplicator(name string) *ChannelDuplicator {
	chDoup := &ChannelDuplicator{
		lock:      lock.NewLock("channel-duplicator"),
		outputs:   make([]chan interface{}, 0),
		inputs:    make([]<-chan interface{}, 0),
		transfer:  make(chan interface{}, 100),
		debug:     false,
		copy:      true,
		debugName: name,
	}
	chDoup.startDuplicator()

	return chDoup
}
func (ch *ChannelDuplicator) EnableCopyMode() {
	ch.copy = true
}

func (ch *ChannelDuplicator) EnableDebug() {
	ch.debug = true
}

func (ch *ChannelDuplicator) SetName(name string) {
	ch.debugName = name
}

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

func (ch *ChannelDuplicator) GetBufferedOutput(buffSize int64) chan interface{} {
	// make a channel with a 10 buffer size
	if ch.debug {
		fmt.Println("adding output on", ch.debugName)
	}
	newOutput := make(chan interface{}, buffSize)
	ch.outputs = append(ch.outputs, newOutput)
	return newOutput
}

func (ch *ChannelDuplicator) UnregisterOutput(remove chan interface{}) {
	var removeIndex int
	for i, channel := range ch.outputs {
		if channel == remove {
			removeIndex = i
		}
	}
	//Remove channel by swapping the removed channel to the end and then just trimming the slice
	ch.outputs[len(ch.outputs)-1], ch.outputs[removeIndex] = ch.outputs[removeIndex], ch.outputs[len(ch.outputs)-1]
	ch.outputs = ch.outputs[:len(ch.outputs)-1]
}

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

func (ch *ChannelDuplicator) Offer(value interface{}) {
	if ch.debug {
		fmt.Println("offering to transfer", ch.debugName)
	}
	if ch.copy && reflect.TypeOf(value).Kind() == reflect.Ptr {
		newVal := deepcopy.Copy(value)
		//pass that pointer down the transfer line
		if ch.debug {
			str, _ := json.Marshal(newVal)
			fmt.Println("offering copy to trasfer=", ch.debugName, "value=", string(str))
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

func (ch *ChannelDuplicator) startDuplicator() {
	go func() {
		for nextValue := range ch.transfer {
			ch.lock.Acquire("startDuplicator")

			if ch.debug {
				fmt.Println("sending down outputs on", ch.debugName)
			}
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
			ch.lock.Release()
		}
	}()

}

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
