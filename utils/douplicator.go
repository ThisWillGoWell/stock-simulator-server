package utils

import (
	"fmt"
	"time"
)

type ChannelDuplicator struct {
	inputs   []chan interface{}
	transfer chan interface{}
	outputs  []chan interface{}
}

func MakeDuplicator() *ChannelDuplicator {
	chDoup := &ChannelDuplicator{
		outputs:  make([]chan interface{}, 0),
		transfer: make(chan interface{}, 10),
	}

	chDoup.startDuplicator()

	return chDoup
}

func (ch *ChannelDuplicator) RegisterOutput() chan interface{} {
	// make a channel with a 10 buffer size
	newChannel := make(chan interface{}, 10)
	ch.outputs = append(ch.outputs, newChannel)
	return newChannel
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

func (ch *ChannelDuplicator) RegisterInput(inputChannel chan interface{}) {
	go func() {
		for val := range inputChannel {
			fmt.Println("got from input", val)
			ch.transfer <- val
		}

	}()
	fmt.Println("returning")
}

func (ch *ChannelDuplicator) Offer(value interface{}) {
	select {
	case ch.transfer <- value:
	default:
	}
}

func (ch *ChannelDuplicator) startDuplicator() {
	go func() {
		for nextValue := range ch.transfer {
			fmt.Println("got from transfer", nextValue)
			for _, channel := range ch.outputs {
				select {
				case channel <- nextValue:
					fmt.Println("sent to output")
				default:
					fmt.Println("missed value")
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
		output := chDoup.RegisterOutput()
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
