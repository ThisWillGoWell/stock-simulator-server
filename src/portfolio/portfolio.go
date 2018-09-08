package portfolio

import (
	"errors"
	"fmt"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
)

const (
	ObjectType = "portfolio"
)

var Portfolios = make(map[string]*Portfolio)
var PortfoliosLock = lock.NewLock("portfolios")
var UpdateChannel = duplicator.MakeDuplicator("portfolio-update")
var NewObjectChannel = duplicator.MakeDuplicator("new-portfolio")

/**
Portfolios are the $$$ part of a user
*/
type Portfolio struct {
	UserUUID string  `json:"user_uuid"`
	UUID     string  `json:"uuid"`
	Wallet   int64 `json:"wallet" change:"-"`
	NetWorth int64 `json:"net_worth" change:"-"`

	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	// stock_uuid -> ledgerObject

	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
	UpdateInput   *duplicator.ChannelDuplicator `json:"-"`

	Lock *lock.Lock `json:"-"`
}

func (port *Portfolio) GetId() string {
	return port.UUID
}

func (port *Portfolio) GetType() string {
	return ObjectType
}

func NewPortfolio(userUUID string) (*Portfolio, error) {
	uuid := utils.PseudoUuid()
	return MakePortfolio(uuid, userUUID, 1000)
}

func MakePortfolio(uuid, userUUID string, wallet int64) (*Portfolio, error) {
	//PortfoliosUpdateChannel.EnableDebug("port update")
	PortfoliosLock.Acquire("new-portfolio")
	defer PortfoliosLock.Release()
	if _, exists := Portfolios[uuid]; exists {
		utils.RemoveUuid(uuid)
		return nil, errors.New("portfolio uuid already Exists")
	}
	port :=
		&Portfolio{
			UserUUID:      userUUID,
			UUID:          uuid,
			Wallet:        wallet,
			UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-update", uuid)),
			Lock:          lock.NewLock(fmt.Sprintf("portfolio-%s", uuid)),
			UpdateInput:   duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-valueable-update", uuid)),
			NetWorth: wallet,
		}
	Portfolios[uuid] = port
	port.Lock.EnableDebug()
	port.UpdateChannel.EnableCopyMode()
	NewObjectChannel.Offer(port)
	UpdateChannel.RegisterInput(port.UpdateChannel.GetOutput())
	utils.RegisterUuid(uuid, port)
	go port.valuableUpdate()
	return port, nil
}

/**
async code that gets called whenever a stock or a ledger that the portfolio owns changes
This then triggers a recalc of net worth and offers its self up as a update
*/
func (port *Portfolio) valuableUpdate() {
	updateChannel := port.UpdateInput.GetOutput()
	for range updateChannel {
		port.Update()
	}
}

func (port *Portfolio) Update(){
	// need to acquire here or else deadlock on the trade
	ledger.EntriesLock.Acquire("portfolio-update")
	defer ledger.EntriesLock.Release()
	port.Lock.Acquire("portfolio-update")
	newNetWorth := port.calculateNetWorth()
	port.NetWorth = newNetWorth
	port.UpdateChannel.Offer(port)
	port.Lock.Release()
}

func GetPortfolio(userUUID string) (*Portfolio, error) {
	port, exists := Portfolios[userUUID]
	if !exists {
		return nil, errors.New("uuid does not have a portfolio tied to it")
	}
	return port, nil
}

//update the current net worth. NOT THREAD SAFE
func (port *Portfolio) calculateNetWorth() int64 {

	sum := int64(0)
	for valueStr, entry := range ledger.EntriesPortfolioStock[port.UUID] {
		value := valuable.Stocks[valueStr]
		sum += value.GetValue() * entry.Amount
	}
	return sum + port.Wallet
}

func (port *Portfolio) RegisterUpdate(newInput chan interface{}) {

}

func GetAllPortfolios() []*Portfolio {
	PortfoliosLock.Acquire("get all ports")
	defer PortfoliosLock.Release()
	lst := make([]*Portfolio, len(Portfolios))
	i := 0
	for _, val := range Portfolios {
		lst[i] = val
		i += 1
	}
	return lst
}
