package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stock-simulator-server/src/items"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/histroy"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/notification"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
)

var clients = make(map[*Client]bool)
var clientsLock = lock.NewLock("clients-lock")
var connections = make(map[string]map[int]*Client)
var currentId = 0

//var BroadcastMessages = duplicator.MakeDuplicator("client-broadcast-messages")
//var currentId = 1
//var Updates = duplicator.MakeDuplicator("Client Updates")
//var NewObjects = duplicator.MakeDuplicator("Client New Objects")

/**
A client is a individual connection
talks directly with the socket connections
A client is tied to a user, there can be many clients to a single user
Though messages that are responded to (query/trades) are only sent back
to the client that sent it
*/
type Client struct {
	socketRx    chan string
	socketTx    chan string
	clientNum   int
	close       chan interface{}
	broadcastTx chan messages.Message
	broadcastRx chan messages.Message

	ws websocket.Conn

	user   *account.User
	active bool
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
	if initialMessage.IsConnect() {
		userTemp, err := account.ConnectUser(initialMessage.Msg.(*messages.ConnectMessage).SessionToken)
		if err != nil {
			return err
		}
		user = userTemp
		sessionToken = initialMessage.Msg.(*messages.ConnectMessage).SessionToken
	} else {
		return errors.New("unknown message, need sessio")
	}

	client := &Client{
		clientNum: currentId,
		user:      user,
		socketRx:  rx,
		socketTx:  tx,
		active:    true,
		close:     make(chan interface{}),
	}
	currentId += 1
	_, exists := connections[user.Uuid]
	if !exists {
		connections[user.Uuid] = make(map[int]*Client)
	}
	connections[user.Uuid][client.clientNum] = client

	go client.tx(sessionToken)
	go client.rx()
	return nil
}

/**
When a session is started, loop though all current cache and send them to the client
Also send the success login to make sure that happens first on the login
*/
func (client *Client) tx(sessionToken string) {
	client.sendMessage(messages.SuccessConnect(client.user.Uuid, sessionToken, client.user.Config))

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
	for _, v := range notification.GetAllNotifications(client.user.Uuid) {
		client.sendMessage(messages.BuildNotificationMessage(v))
	}
	for _, v := range items.GetItemsForUser(client.user.Uuid) {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	//finally register the broadcast message as a input to the clients message sender
	//do this after the payload to prevent an update from being sent before the object prototype is sent
	send := client.user.Sender.GetOutput()
	go func() {
		for msg := range send {
			client.sendMessage(msg)
		}
	}()
	<-client.close
	client.user.Sender.CloseOutput(send)
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
			client.sendMessage(messages.NewErrorMessage("err unmarshaling json"))
			continue
		}

		switch message.Action {
		case messages.NotificationAck:
			client.processAckMessage(message)
		case messages.ChatAction:
			client.processChatMessage(message.Msg.(messages.Message))
		case messages.TradeAction:
			client.processTradeMessage(message)
		case messages.QueryAction:
			client.processQueryMessage(message)
		case messages.TransferAction:
			client.processTransferMessage(message)
		case messages.SetAction:
			client.processSetMessage(message)
		default:
			client.sendMessage(messages.NewErrorMessage("action is not known"))
		}
	}
	client.active = false
	clientsLock.Acquire("remove uuid from connections")
	defer clientsLock.Release()
	delete(connections[client.user.Uuid], client.clientNum)
	if len(connections[client.user.Uuid]) == 0 {
		delete(connections, client.user.Uuid)
	}
	client.close <- nil
	client.user.LogoutUser()
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
	account.GlobalMessages.Offer(messages.BuildChatMessage(message.(*messages.ChatMessage)))
}

func (client *Client) processTradeMessage(baseMessage *messages.BaseMessage) {
	tradeMessage := baseMessage.Msg.(*messages.TradeMessage)
	po := order.MakePurchaseOrder(tradeMessage.StockId, client.user.PortfolioId, tradeMessage.Amount)
	go func() {
		response := <-po.ResponseChannel
		client.user.Sender.Output.Offer(messages.BuildResponseMsg(response, baseMessage.RequestID))
	}()
}

func (client *Client) processTransferMessage(baseMessage *messages.BaseMessage) {

	transferMessage := baseMessage.Msg.(*messages.TransferMessage)
	to := order.MakeTransferOrder(client.user.PortfolioId, transferMessage.Recipient, transferMessage.Amount)
	go func() {
		response := <-to.ResponseChannel
		client.sendMessage(messages.BuildResponseMsg(response, baseMessage.RequestID))

	}()
}

func (client *Client) processAckMessage(baseMessage *messages.BaseMessage) {
	ackMessage := baseMessage.Msg.(*messages.NotificationAckMessage)
	err := notification.AcknowledgeNotification(ackMessage.Uuid, client.user.Uuid)
	if err != nil {
		client.sendMessage(messages.NewErrorMessage(err.Error()))
	}
}

func (client *Client) processQueryMessage(baseMessage *messages.BaseMessage) {
	queryMessage := baseMessage.Msg.(*messages.QueryMessage)
	q := histroy.MakeQuery(queryMessage)
	go func() {
		response := <-q.ResponseChannel
		client.sendMessage(messages.BuildResponseMsg(response, baseMessage.RequestID))
	}()
}

func (client *Client) processItemMessage(m messages.BaseMessage) {
	itemMessage := m.Msg.(messages.ItemMessage)
	switch itemMessage.O.(type) {
	case messages.ItemBuyMessage:
		err := items.BuyItem(client.user.Uuid, itemMessage.O.(messages.ItemBuyMessage).ItemName)
		if err != nil {
			client.sendMessage(messages.BuildItemBuyFailedMessage(itemMessage.O.(messages.ItemBuyMessage).ItemName, m.RequestID, err))
		}
	case messages.ItemViewMessage:
		result, err := items.ViewItem(itemMessage.O.(messages.ItemViewMessage).ItemUuid, client.user.Uuid)
		if err != nil {
			client.sendMessage(messages.BuildItemViewFailedMessage(itemMessage.O.(messages.ItemViewMessage).ItemUuid, m.RequestID, err))
		} else {
			client.sendMessage(messages.BuildItemViewMessage(itemMessage.O.(messages.ItemViewMessage).ItemUuid, m.RequestID, result))
		}
	case messages.ItemUseMessage:
		result, err := items.Use(itemMessage.O.(messages.ItemUseMessage).ItemUuid, client.user.Uuid, itemMessage.O.(messages.ItemUseMessage).UseParameters)
		if err != nil {
			client.sendMessage(messages.BuildItemUseFailedMessage(itemMessage.O.(messages.ItemUseMessage).ItemUuid, m.RequestID, err))
		} else {
			client.sendMessage(messages.BuildItemUseMessage(itemMessage.O.(messages.ItemUseMessage).ItemUuid, m.RequestID, result))
		}
	}
}

func (client *Client) processSetMessage(baseMessage *messages.BaseMessage) {
	setMessage := baseMessage.Msg.(*messages.SetMessage)
	var response *messages.SetResponse
	switch setMessage.Set {
	case "config":
		newConfig, ok := setMessage.Value.(map[string]interface{})
		if !ok {
			response = messages.BuildFailedSet(errors.New("failed, invalid type"))
			break
		}
		client.user.SetConfig(newConfig)
		response = messages.BuildSuccessSet()
	case "password":
		newPassword, ok := setMessage.Value.(string)
		if !ok {
			response = messages.BuildFailedSet(errors.New("failed, invalid type"))
			break
		}
		err := client.user.SetPassword(newPassword)
		if err != nil {
			response = messages.BuildFailedSet(err)
			break
		}
		response = messages.BuildSuccessSet()
	case "display_name":
		newDisplayName, ok := setMessage.Value.(string)
		if !ok {
			response = messages.BuildFailedSet(errors.New("failed, invalid type"))
			break
		}
		err := client.user.SetDisplayName(newDisplayName)
		if err != nil {
			response = messages.BuildFailedSet(err)
			break
		}
		response = messages.BuildSuccessSet()
	}
	client.sendMessage(messages.BuildResponseMsg(response, baseMessage.RequestID))

}
