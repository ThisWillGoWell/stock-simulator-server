package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/trade"
	"github.com/stock-simulator-server/src/valuable"
)

var clients = make(map[*Client]bool)
var clientsLock = lock.NewLock("clients-lock")

var BroadcastMessages = duplicator.MakeDuplicator("client-broadcast-messages")

func BroadcastMessageBuilder() {
	updates := change.SubscribeUpdateOutput.GetBufferedOutput(100)
	go func() {
		for update := range updates {
			BroadcastMessages.Offer(messages.BuildUpdateMessage(update))
		}
	}()
	//BroadcastMessagePrinter()
}

func BroadcastMessagePrinter() {

	msgs := BroadcastMessages.GetBufferedOutput(100)
	go func() {
		for msg := range msgs {
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

	messageSender *duplicator.ChannelDuplicator

	broadcastTx chan messages.Message
	broadcastRx chan messages.Message

	ws websocket.Conn

	user   *account.User
	active bool
}

func InitialRecieve(initialPayload string, tx, rx chan string) error {
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
		user = userTemp
	}

	client := &Client{
		user:          user,
		socketRx:      rx,
		socketTx:      tx,
		messageSender: duplicator.MakeDuplicator("client-" + user.Uuid + "-message"),
		active:        true,
	}
	client.tx()
	go client.rx()
	go client.initSession()
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
			client.initSession()
		default:
			client.messageSender.Offer(messages.NewErrorMessage("action is not known"))
		}
	}
	client.active = false
	client.user.LogoutUser()
}

// send down websocket
func (client *Client) tx() {
	send := client.messageSender.GetBufferedOutput(100)
	go func() {
		for msg := range send {
			client.sendMessage(msg)
		}
	}()
}

func (client *Client) sendMessage(msg interface{}) {
	str, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	client.socketTx <- string(str)
}

func (client *Client) processChatMessage(message messages.Message) {
	chatMessage := message.(*messages.ChatMessage)
	chatMessage.Author = client.user.Uuid
	chatMessage.Timestamp = time.Now().Unix()
	BroadcastMessages.Offer(chatMessage)
}

func (client *Client) processTradeMessage(message messages.Message) {
	tradeMessage := message.(*messages.TradeMessage)
	po := order.BuildPurchaseOrder(tradeMessage.StockId, tradeMessage.ExchangeID, client.user.Uuid, tradeMessage.Amount)
	trade.Trade(po)
	go func() {
		response := <-po.ResponseChannel
		client.sendMessage(messages.BuildPurchaseResponse(response))
	}()
}

func (client *Client) initSession() {

	client.sendMessage(messages.SuccessLogin(client.user.Uuid))
	for _, v := range account.GetAllUsers() {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	for _, v := range portfolio.GetAllPortfolios() {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	for _, v := range valuable.GetAllStocks() {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	for _, v := range ledger.GetAllLedgers() {
		client.sendMessage(messages.NewObjectMessage(v))
	}

	client.messageSender.RegisterInput(BroadcastMessages.GetBufferedOutput(50))
}
