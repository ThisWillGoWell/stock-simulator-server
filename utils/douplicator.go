package utils

import (
	"fmt"
	"time"
)

type ChannelDuplicator struct {
	transfer chan interface{}
	outputs  []chan interface{}
}

func MakeDuplicator() *ChannelDuplicator {
	chDoup := &ChannelDuplicator{
		outputs:  make([]chan interface{}, 100),
		transfer: make(chan interface{}, 100),
	}

	chDoup.startDuplicator()

	return chDoup
}

func (ch *ChannelDuplicator) GetOutput() chan interface{} {
	// make a channel with a 10 buffer size
	newOutput := make(chan interface{}, 10)
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

func (ch *ChannelDuplicator) RegisterInput(inputChannel <- chan interface{}) {
	go func() {
		for val := range inputChannel {
			ch.transfer <- val
		}
	}()

}

func (ch *ChannelDuplicator) Offer(value interface{}) {
	ch.transfer <- value
}

func (ch *ChannelDuplicator) startDuplicator() {
	go func() {
		for nextValue := range ch.transfer {
			for _, channel := range ch.outputs {
				select {
				case channel <- nextValue:
				default:
				}
			}
		}
	}()

}

func main() {
	input1 := make(chan interface{})
	input2 := make(chan interface{})

	chDoup := MakeDuplicator()
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
