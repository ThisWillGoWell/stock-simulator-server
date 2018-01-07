package client

import (
	"stock-server/utils"
	"stock-server/wallstreet"
	"github.com/gorilla/websocket"
	"stock-server/messages"
)

var clients = make(map[*Client]bool)
var clientsLock = utils.NewLock()
var clientBroadcast = utils.MakeDuplicator()

type Client struct {
	socketRx chan *messages.BaseMessage
	socketTx chan *messages.BaseMessage

	broadcastTx chan *messages.BaseMessage
	broadcastRx chan *messages.BaseMessage


	stockUpdate     chan interface{}
	portfolioUpdate chan interface{}
	ws              websocket.Conn

	portfolio *wallstreet.Portfolio
	user *User
}


func Login(username, password string, tx, rx chan *messages.BaseMessage) (error){
	user, err := getUser(username, password)
	if err != nil {
		return err
	}
	portfilio := wallstreet.PortfolioIds[user.uuid]

	client := &Client {
		socketRx:        rx,
		socketTx:        tx,
		stockUpdate:     portfilio.Exchange.GetStockUpdateChanel(),
		portfolioUpdate: portfilio.UpdateChannel.GetOutput(),
	}
	go client.rx()
	go client.tx()

	return nil
}



func (client *Client)rx(){
	for message := range client.socketRx {
		if message.IsChat(){
			client.broadcastTx <- message
		}else if message.IsTrade(){
		}
	}
}
func (client *Client) tx(){
	for {
		select{
		case msg := <-client.broadcastRx:
			client.socketTx <- msg
		case stock := <-client.stockUpdate:
			client.socketTx <- messages.NewUpdateMessage(messages.StockUpdate, stock)
		case portfolio := <- client.portfolioUpdate:
			client.socketTx <- messages.NewUpdateMessage(messages.PortfolioUpdate, portfolio)


		}
	}
}

func makeStockUpdate


