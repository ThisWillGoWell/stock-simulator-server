package client

import (
	"github.com/stock-simulator-server/account"
	"github.com/stock-simulator-server/utils"
	"github.com/stock-simulator-server/exchange"
	"github.com/gorilla/websocket"
	"github.com/stock-simulator-server/messages"
	"github.com/stock-simulator-server/portfolio"
	"github.com/stock-simulator-server/valuable"
	"encoding/json"
	"time"
	"github.com/stock-simulator-server/order"
	"errors"
)

var clients = make(map[*Client]bool)
var clientsLock = utils.NewLock("clients-lock")
var clientBroadcast = utils.MakeDuplicator()

var broadcastMessages = utils.MakeDuplicator()


func BroadcastMessageBuilder(){
	valuableUpdateChannel := valuable.ValuableUpdateChannel.GetOutput()
	portfolioUpdateChannel := portfolio.PortfoliosUpdateChannel.GetOutput()
	ledgerUpdateChannel := portfolio.PortfoliosUpdateChannel.GetOutput()
	for{
		var update interface{}
		var updateType string
		select{
		case update = <- portfolioUpdateChannel:
			updateType = messages.PortfolioUpdate
		case update = <- valuableUpdateChannel:
			updateType = messages.ValuableUpdate
		case  update = <- ledgerUpdateChannel:
			updateType = messages.LedgerUpdate
		}
		msg :=messages.BuildUpdateMessage(updateType, update)
		broadcastMessages.Offer(msg)
	}
}


type Client struct {
	socketRx chan string
	socketTx chan string

	messageSender *utils.ChannelDuplicator

	broadcastTx chan messages.Message
	broadcastRx chan messages.Message

	ws              websocket.Conn

	user *account.User
	}


func Login(loginMessageStr string, tx, rx chan string) (error){
	loginBaseMessage := new(messages.BaseMessage)
	unmarshErr := loginBaseMessage.UnmarshalJSON([]byte(loginMessageStr))
	if unmarshErr != nil{
		return unmarshErr
	}
	if !loginBaseMessage.IsLogin(){
		return errors.New("wrong type")
	}
	loginMessage := loginBaseMessage.Value.(*messages.LoginMessage)

	user, err := account.GetUser(loginMessage.Username, loginMessage.Password)
	if err != nil {
		return err
	}

	if err != nil{
		return err
	}

	client := &Client {
		user: user,
		socketRx:        rx,
		socketTx:        tx,
		messageSender: utils.MakeDuplicator(),
		}
	go client.rx()
	go client.tx()
	client.messageSender.RegisterInput(broadcastMessages.GetOutput())

	return nil
}


//receive go routine
func (client *Client)rx(){
	for messageString := range client.socketRx {
		message := new(messages.BaseMessage)
		//attempt to
		err := message.UnmarshalJSON([]byte(messageString))
		if err != nil{
			client.messageSender.Offer(messages.NewErrorMessage(err))
			continue
		}

		switch message.Action {
		case messages.ChatAction:
			client.processChatMessage(message.Value)
		case messages.TradeAction:
			client.processTradeMessage(message.Value)
		case messages.UpdateAction:
			client.processUpdateMessage()
		default:
			client.messageSender.Offer(messages.NewErrorMessage(errors.New("action is not known")))
		}
	}
}
// send down websocket
func (client *Client) tx(){
	send := client.messageSender.GetOutput()
	for sendMsg := range send{
		str, err := json.Marshal(sendMsg)
		if err != nil{
			panic(err)
		}else{
			client.socketTx <- string(str)
		}
	}
}

func (client *Client)processChatMessage(message messages.Message){
	chatMessage := message.(*messages.ChatMessage)
	chatMessage.Author = client.user.Uuid
	chatMessage.Timestamp = time.Now().Unix()
	client.messageSender.Offer(chatMessage)
}

func (client *Client)processTradeMessage(message messages.Message){
	tradeMessage := message.(*messages.TradeMessage)
	po := order.BuildPurchaseOrder(tradeMessage.StockTicker, tradeMessage.ExchangeID, client.user.Uuid, tradeMessage.Amount)
	exchange.InitiateTrade(po)
	go func(){
		response := <- po.ResponseChannel
		client.messageSender.Offer(messages.BuildPurchaseResponse(tradeMessage, response))
	}()
}

func (client *Client)processUpdateMessage() {
	for _, entry := range exchange.Exchanges{
		message := messages.BuildUpdateMessage(messages.LedgerUpdate, entry)
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
}

