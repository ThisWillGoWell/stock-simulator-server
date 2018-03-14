package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/exchange"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/utils"
)

var clients = make(map[*Client]bool)
var clientsLock = utils.NewLock("clients-lock")
var clientBroadcast = utils.MakeDuplicator()

var BroadcastMessages = utils.MakeDuplicator()

func BroadcastMessageBuilder() {
	updates := utils.SubscribeUpdateOutput.GetOutput()
	go func() {
		for update := range updates {
			BroadcastMessages.Offer(messages.BuildUpdateMessage(update))
		}
	}()
	BroadcastMessagePrinter()
}

func BroadcastMessagePrinter() {
	messages := BroadcastMessages.GetOutput()
	return
	go func() {
		for msg := range messages {
			str, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			} else {
				fmt.Println(string(str))
			}
		}
	}()
}

type Client struct {
	socketRx chan string
	socketTx chan string

	messageSender *utils.ChannelDuplicator

	broadcastTx chan messages.Message
	broadcastRx chan messages.Message

	ws websocket.Conn

	user *account.User
	active bool
}

func InitialRecieve(initialPayload string,  tx, rx chan string) error{
	initialMessage := new(messages.BaseMessage)
	unmarshalErr := initialMessage.UnmarshalJSON([]byte(initialPayload))

	if unmarshalErr != nil {
		return unmarshalErr
	}
	user := new(account.User)
	if initialMessage.IsAccountCreate() {
		userTemp, err := account.NewUser(initialMessage.Msg.(*messages.NewAccountMessage).UserName, initialMessage.Msg.(*messages.NewAccountMessage).Password)
		if err != nil {
			return err
		}
		user = userTemp
	} else if initialMessage.IsLogin() {
		userTemp, err := account.GetUser(initialMessage.Msg.(*messages.LoginMessage).UserName, initialMessage.Msg.(*messages.LoginMessage).Password)
		if err != nil {
			return err
		}
		user= userTemp
	}

	client := &Client{
		user:          user,
		socketRx:      rx,
		socketTx:      tx,
		messageSender: utils.MakeDuplicator(),
		active: true,
	}
	client.tx()
	go client.rx()
	client.messageSender.RegisterInput(BroadcastMessages.GetBufferedOutput(50))
	client.sendInitialPayload()
	return nil

}


//receive go routine
func (client *Client) rx() {

	for messageString := range client.socketRx {
		message := new(messages.BaseMessage)
		//attempt to
		err := message.UnmarshalJSON([]byte(messageString))
		if err != nil {
			client.messageSender.Offer(messages.NewErrorMessage("err unmarshaling json"))
			continue
		}

		switch message.Action {
		case messages.ChatAction:
			client.processChatMessage(message.Msg.(messages.Message))
		case messages.TradeAction:
			client.processTradeMessage(message.Msg.(messages.Message))
		case messages.UpdateAction:
			client.sendInitialPayload()
		default:
			client.messageSender.Offer(messages.NewErrorMessage("action is not known"))
		}
	}
	client.active = false
	client.user.LogoutUser()
}

// send down websocket
func (client *Client) tx() {
	send := client.messageSender.GetOutput()
	batchSendTicker := time.NewTicker(1 * time.Second)
	sendQueue := make(chan interface{}, 300)
	go func(){
		for {
			select {
			case <-batchSendTicker.C:
				if ! client.active{
					break
				}
				sendOutQueue(sendQueue, client.socketTx)
			case msg := <-send:
				select {
				case sendQueue <- msg:
				default:
					//the queue is full
					//empty it
					sendOutQueue(sendQueue, client.socketTx)
					//add
					sendQueue <- msg
				}

			}
		}
	}()
	//

}
func sendOutQueue(sendQueue chan interface{}, socketTx chan string) {
	sendList := make([]interface{}, 0)
emptyQueue:
	for {
		select {
		case ele := <-sendQueue:
			sendList = append(sendList, ele)
		default:
			break emptyQueue
		}
	}

	if len(sendList) > 0 {
		str, err := json.Marshal(sendList)
		if err != nil {
			panic(err)
		} else {
			socketTx <- string(str)
		}
	}
}

func (client *Client) processChatMessage(message messages.Message) {
	chatMessage := message.(*messages.ChatMessage)
	chatMessage.Author = client.user.Uuid
	chatMessage.Timestamp = time.Now().Unix()
	BroadcastMessages.Offer(chatMessage)
}

func (client *Client) processTradeMessage(message messages.Message) {
	tradeMessage := message.(*messages.TradeMessage)
	po := order.BuildPurchaseOrder(tradeMessage.StockTicker, tradeMessage.ExchangeID, client.user.Uuid, tradeMessage.Amount)
	exchange.InitiateTrade(po)
	go func() {
		response := <-po.ResponseChannel
		client.messageSender.Offer(messages.BuildPurchaseResponse(tradeMessage, response))
	}()
}

func (client *Client) sendInitialPayload() {
	fmt.Println("got update")
	client.messageSender.Offer(messages.BuildUpdateMessage(utils.GetCurrentValues()))
	/*
		for _, entry := range exchange.Exchanges{
			message := messages.BuildUpdateMessage(messages.LedgerUpdate, entry.Ledger)
			client.messageSender.Offer(message)
		}

		for _, entry := range portfolio.Portfolios{
			message := messages.BuildUpdateMessage(messages.PortfolioUpdate, entry)
			client.messageSender.Offer(message)
		}

		for _, entry := range valuable.Valuables{
			message := messages.BuildUpdateMessage(messages.ValuableUpdate, entry)
			client.messageSender.Offer(message)
		}
	*/
}
