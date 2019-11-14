package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"

	metrics "github.com/ThisWillGoWell/stock-simulator-server/src/app/metics"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/record"

	"github.com/ThisWillGoWell/stock-simulator-server/src/messages"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/items"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/sender"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database/histroy"
	"github.com/ThisWillGoWell/stock-simulator-server/src/game/order"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/ledger"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/user"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

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

	user   *user.User
	active bool
}

/**
Because for some reason, the js web socket class does not have headers,
We accept all connections and pull the initial payload as a login/user create command
*/
func InitialReceive(initialPayload string, tx, rx chan string) error {
	log.Log.Info("initial recieve of new client", initialPayload)
	clientsLock.Acquire("initial received of new client")
	defer clientsLock.Release()
	initialMessage := new(messages.BaseMessage)
	unmarshalErr := initialMessage.UnmarshalJSON([]byte(initialPayload))

	if unmarshalErr != nil {
		log.Log.Error("Unmarshal error: ", initialPayload)
		return unmarshalErr
	}
	u := new(user.User)
	var sessionToken string
	if initialMessage.IsConnect() {
		userTemp, err := user.ConnectUser(initialMessage.Msg.(*messages.ConnectMessage).SessionToken)
		if err != nil {
			log.Log.Error("error in connecting user: ", err, u.Uuid)
			return err
		}
		u = userTemp
		sessionToken = initialMessage.Msg.(*messages.ConnectMessage).SessionToken
	} else {
		log.Log.Error("unknown message action for connect message", initialPayload)
		return errors.New("unknown message, need session")
	}

	client := &Client{
		clientNum: currentId,
		user:      u,
		socketRx:  rx,
		socketTx:  tx,
		active:    true,
		close:     make(chan interface{}),
	}
	currentId += 1
	_, exists := connections[u.Uuid]
	if !exists {
		connections[u.Uuid] = make(map[int]*Client)
	}
	connections[u.Uuid][client.clientNum] = client
	log.Log.Info("client connected ", u.Uuid, u.ActiveClients)
	metrics.ClientConnect()
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
	for _, v := range user.GetAllUsers() {
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
	for _, v := range notification.GetAllNotifications(client.user.PortfolioId) {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	for _, v := range items.GetItemsForPortfolio(client.user.PortfolioId) {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	books, records := record.GetRecordsForPortfolio(client.user.PortfolioId)
	for _, v := range books {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	for _, v := range records {
		client.sendMessage(messages.NewObjectMessage(v))
	}
	for _, e := range effect.GetAllEffects() {
		client.sendMessage(messages.NewObjectMessage(e))
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
		metrics.RecieveMessage(len(messageString))
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
		case messages.ItemAction:
			client.processItemMessage(message)
		case messages.LevelUpAction:
			client.processLevelUpAction(message)
		case messages.DeleteAction:
			client.processDeleteAction(message)
		case messages.ProspectTradeAction:
			client.processProspectMessage(message)
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
	metrics.ClientDisconnect()
}

/**
blocking send a single message
*/
func (client *Client) sendMessage(msg interface{}) {
	str, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	s := string(str)
	metrics.SendMessage(len(s))
	client.socketTx <- s
}

func (client *Client) processChatMessage(message messages.Message) {
	chatMessage := message.(*messages.ChatMessage)
	chatMessage.Author = client.user.Uuid
	chatMessage.Timestamp = time.Now()
	//database.SaveChatMessage(chatMessage.Author, chatMessage.Message)
	sender.GlobalMessages.Offer(messages.BuildChatMessage(message.(*messages.ChatMessage)))
}

func (client *Client) processTradeMessage(baseMessage *messages.BaseMessage) {
	tradeMessage := baseMessage.Msg.(*messages.TradeMessage)
	po := order.MakePurchaseOrder(tradeMessage.StockId, client.user.PortfolioId, tradeMessage.Amount)
	go func() {
		response := <-po.ResponseChannel
		client.sendMessage(messages.BuildResponseMsg(response, baseMessage.RequestID))
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
	err := notification.AcknowledgeNotification(ackMessage.Uuid, client.user.PortfolioId)
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

func (client *Client) processItemMessage(m *messages.BaseMessage) {
	itemMessage := m.Msg.(*messages.ItemMessage)
	switch itemMessage.O.(type) {
	case *messages.ItemBuyMessage:
		uuid, err := items.BuyItem(client.user.PortfolioId, itemMessage.O.(*messages.ItemBuyMessage).ItemConfig)
		if err != nil {
			client.sendMessage(messages.BuildItemBuyFailedMessage(itemMessage.O.(*messages.ItemBuyMessage).ItemConfig, m.RequestID, err))
		} else {
			client.sendMessage(messages.BuildItemBuySuccessMessage(itemMessage.O.(*messages.ItemBuyMessage).ItemConfig, m.RequestID, uuid))
		}
	//case *messages.ItemViewMessage:
	//	result, err := items.ViewItem(itemMessage.O.(*messages.ItemViewMessage).ItemUuid, client.user.Uuid)
	//	if err != nil {
	//		client.sendMessage(messages.BuildItemViewFailedMessage(itemMessage.O.(*messages.ItemViewMessage).ItemUuid, m.RequestID, err))
	//	} else {
	//		client.sendMessage(messages.BuildItemViewMessage(itemMessage.O.(*messages.ItemViewMessage).ItemUuid, m.RequestID, result))
	//	}
	case *messages.ItemUseMessage:
		result, err := items.Use(itemMessage.O.(*messages.ItemUseMessage).ItemUuid, client.user.PortfolioId, itemMessage.O.(*messages.ItemUseMessage).UseParameters)
		if err != nil {
			client.sendMessage(messages.BuildItemUseFailedMessage(itemMessage.O.(*messages.ItemUseMessage).ItemUuid, m.RequestID, err))
		} else {
			client.sendMessage(messages.BuildItemUseMessage(itemMessage.O.(*messages.ItemUseMessage).ItemUuid, m.RequestID, result))
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
		if err := client.user.SetConfig(newConfig); err != nil {
			response = messages.BuildFailedSet(err)
		} else {
			response = messages.BuildSuccessSet()
		}

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

func (client *Client) processLevelUpAction(baseMessage *messages.BaseMessage) {
	err := portfolio.Portfolios[client.user.PortfolioId].LevelUp()
	client.sendMessage(messages.BuildLevelUpResponse(baseMessage.RequestID, err))
}

func (client *Client) processDeleteAction(baseMessage *messages.BaseMessage) {
	deleteMsg := baseMessage.Msg.(*messages.DeleteMessage)
	var err error
	switch deleteMsg.Type {
	case objects.ItemIdentifiableType:
		err = items.DeleteRequest(deleteMsg.Uuid)
	case objects.NotificationIdentifiableType:
		err = notification.DeleteRequest(deleteMsg.Uuid)
	}
	client.sendMessage(messages.BuildDeleteResponseMsg(baseMessage.RequestID, err))
}

func (client *Client) processProspectMessage(baseMessage *messages.BaseMessage) {
	prospectMessage := baseMessage.Msg.(*messages.TradeMessage)
	prospect := order.MakeProspect(prospectMessage.StockId, client.user.PortfolioId, prospectMessage.Amount)
	go func() {
		response := <-prospect.ResponseChannel
		client.sendMessage(messages.BuildResponseMsg(response, baseMessage.RequestID))
	}()
}
