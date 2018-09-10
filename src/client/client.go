package client

import (
	"encoding/json"
	"fmt"
	"github.com/stock-simulator-server/src/session"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/histroy"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
)

var clients = make(map[*Client]bool)
var clientsLock = lock.NewLock("clients-lock")
var connections = make(map[string]map[int]*Client)
var BroadcastMessages = duplicator.MakeDuplicator("client-broadcast-messages")
var currentId = 1
var Updates = duplicator.MakeDuplicator("Client Updates")
var NewObjects = duplicator.MakeDuplicator("Client New Objects")

/*
This is responsive for accepting message duplicators and converting those objects
into messages to then be fanned out to all clients
*/
func BroadcastMessageBuilder() {
	//BroadcastMessages.EnableDebug()
	updates := Updates.GetBufferedOutput(100)
	go func() {
		for update := range updates {
			BroadcastMessages.Offer(messages.BuildUpdateMessage(update))
		}
	}()
	newObjects := NewObjects.GetBufferedOutput(100)
	go func() {
		for newObject := range newObjects {
			BroadcastMessages.Offer(messages.NewObjectMessage(newObject.(change.Identifiable)))
		}
	}()
	//BroadcastMessagePrinter()
}

/**
Debugging: prints all messages that are sent down a broadcast
*/
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

/**
A client is a individual connection
talks directly with the socket connections
A client is tied to a user, there can be many clients to a single user
Though messages that are responded to (query/trades) are only sent back
to the client that sent it
*/
type Client struct {
	socketRx chan string
	socketTx chan string
	clientNum int
	messageSender *duplicator.ChannelDuplicator

	broadcastTx chan messages.Message
	broadcastRx chan messages.Message

	ws websocket.Conn

	user   *account.User
	active bool
}

func SendToUser(uuid string, value interface{} ){
	clientsLock.Acquire("send-to-clients")
	defer clientsLock.Release()
	clients, exist := connections[uuid]
	if !exist {
		return
	}
	for i := range clients{
		connections[uuid][i].messageSender.Offer(value)
	}
}

/**
Because for some reason, the js web socket class does not have headers,
We accept all connections and pull the initial payload as a login/account create command
*/
func InitialReceive(initialPayload string, tx, rx chan string) error {
	clientsLock.Acquire("initial received of new client")
	defer clientsLock.Release()
	initialMessage := new(messages.BaseMessage)
	unmarshalErr := initialMessage.UnmarshalJSON([]byte(initialPayload))

	if unmarshalErr != nil {
		return unmarshalErr
	}
	user := new(account.User)
	var sessionToken string

	if initialMessage.IsAccountCreate() {
		userTemp, err := account.NewUser(initialMessage.Msg.(*messages.NewAccountMessage).UserName, initialMessage.Msg.(*messages.NewAccountMessage).Password)
		if err != nil {
			return err
		}
		user = userTemp
		sessionToken = session.NewSessionToken(user.Uuid)
	} else if initialMessage.IsLogin() {
		userTemp, err := account.GetUser(initialMessage.Msg.(*messages.LoginMessage).UserName, initialMessage.Msg.(*messages.LoginMessage).Password)
		if err != nil {
			return err
		}
		user = userTemp
		sessionToken = session.NewSessionToken(user.Uuid)
	} else if initialMessage.IsRenew(){
		userTemp, err := account.RenewUser(initialMessage.Msg.(*messages.RenewMessage).SessionToken)
		if err != nil {
			return err
		}
		user = userTemp
		sessionToken = initialMessage.Msg.(*messages.RenewMessage).SessionToken
	}

	client := &Client{
		clientNum: currentId,
		user:          user,
		socketRx:      rx,
		socketTx:      tx,
		messageSender: duplicator.MakeDuplicator("client-" + user.Uuid + "-message"),
		active:        true,
	}
	currentId += 1
	_, exists := connections[user.Uuid]
	if !exists{
		connections[user.Uuid] = make(map[int]*Client)
	}
	connections[user.Uuid][client.clientNum] = client



	client.tx()
	go client.rx()
	go client.initSession(sessionToken)
	return nil

}

/**
go routine for handling the rx portion of the socket
*/
func (client *Client) rx() {
	for messageString := range client.socketRx {
		fmt.Println("MSG: " + messageString)
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
		case messages.QueryAction:
			client.processQueryMessage(message.Msg.(messages.Message))

		default:
			client.messageSender.Offer(messages.NewErrorMessage("action is not known"))
		}
	}
	client.active = false
	clientsLock.Acquire("remove uuid from connections")
	defer clientsLock.Release()
	delete(connections[client.user.Uuid], client.clientNum)
	if len(connections[client.user.Uuid]) == 0{
		delete(connections, client.user.Uuid)
	}
	client.user.LogoutUser()
}

/**
handle the transmit portion of the socket off the clients message sender douplicator
*/
func (client *Client) tx() {
	// note the buffered output, number 100 was completely arbitrary
	send := client.messageSender.GetBufferedOutput(100)
	go func() {
		for msg := range send {
			client.sendMessage(msg)
		}
	}()
}

/**
blocking send a single message
*/
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
	chatMessage.Timestamp = time.Now()
	//database.SaveChatMessage(chatMessage.Author, chatMessage.Message)
	print("offer")
	BroadcastMessages.Offer(messages.BuildChatMessage(message.(*messages.ChatMessage)))
}

func (client *Client) processTradeMessage(message messages.Message) {
	tradeMessage := message.(*messages.TradeMessage)
	po := order.MakePurchaseOrder(tradeMessage.StockId, client.user.PortfolioId, tradeMessage.Amount)
	go func() {
		response := <-po.ResponseChannel
		SendToUser(client.user.Uuid, messages.BuildPurchaseResponse(response))
	}()
}

func (client *Client) processTransferMessage(message messages.Message) {
	transferMessage := message.(*messages.TransferMessage)
	 to := order.MakeTransferOrder(client.user.PortfolioId, transferMessage.Recipient, transferMessage.Amount)
	go func() {
		response := <-to.ResponseChannel
		SendToUser(client.user.Uuid, messages.BuildTransferResponse(response))
	}()
}

func (client *Client) processQueryMessage(message messages.Message) {
	queryMessage := message.(*messages.QueryMessage)
	q := histroy.MakeQuery(queryMessage)
	go func(){
		client.sendMessage(<- q.ResponseChannel )
	}()
}

/**
When a session is started, loop though all current cache and send them to the client
*/
func (client *Client) initSession(sessionToken string ) {
	client.sendMessage(messages.SuccessLogin(client.user.Uuid, sessionToken))

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
	//finally register the broadcast message as a input to the clients message sender
	//do this after the payload to prevent an update from being sent before the object prototype is sent
	client.messageSender.RegisterInput(BroadcastMessages.GetBufferedOutput(50))
}
