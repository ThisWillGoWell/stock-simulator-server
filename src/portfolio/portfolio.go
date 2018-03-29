package portfolio

import (
	"errors"
	"fmt"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/valuable"
)

const (
	ObjectType = "portfolio"
)

var Portfolios = make(map[string]*Portfolio)
var PortfoliosLock = lock.NewLock("portfolios")
var PortfoliosUpdateChannel = duplicator.MakeDuplicator("portfolio-update")
var NewPortfolioChannel = duplicator.MakeDuplicator("new-portfolio")

type Portfolio struct {
	Name     string  `json:"name"`
	UUID     string  `json:"uuid"`
	Wallet   float64 `json:"wallet" change:"-"`
	NetWorth float64 `json:"net_worth" change:"-"`

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

func NewPortfolio(userUUID, name string) (*Portfolio, error) {
	return MakePortfolio(userUUID, name, 1000)
}

func MakePortfolio(uuid, name string, wallet float64) (*Portfolio, error) {
	//PortfoliosUpdateChannel.EnableDebug("port update")
	PortfoliosLock.Acquire("new-portfolio")
	defer PortfoliosLock.Release()
	if _, exists := Portfolios[uuid]; exists {
		return nil, errors.New("portfolio uuid already Exists")
	}
	port :=
		&Portfolio{
			Name:          name,
			UUID:          uuid,
			Wallet:        wallet,
			UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-update", uuid)),
			Lock:          lock.NewLock(fmt.Sprintf("portfolio-%s", name)),
			UpdateInput:   duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-valueable-update", uuid)),
		}
	Portfolios[uuid] = port
	PortfoliosUpdateChannel.RegisterInput(port.UpdateChannel.GetOutput())
	go port.valuableUpdate()
	//NewPortfolioChannel.Offer(port)
	PortfoliosUpdateChannel.Offer(port)
	return port, nil
}

func (port *Portfolio) valuableUpdate() {
	updateChannel := port.UpdateInput.GetOutput()

	for range updateChannel {
		port.Lock.Acquire("portfolio-update")
		newNetWorth := port.calculateNetWorth()
		port.NetWorth = newNetWorth
		port.UpdateChannel.Offer(port)
		port.Lock.Release()
	}
}

func GetPortfolio(userUUID string) (*Portfolio, error) {
	port, exists := Portfolios[userUUID]
	if !exists {
		return nil, errors.New("uuid does not have a portfolio tied to it")
	}
	return port, nil
}

//update the current net worth. NOT THREAD SAFE
func (port *Portfolio) calculateNetWorth() float64 {
	ledger.EntriesLock.Acquire("calculate-worth")
	defer ledger.EntriesLock.Release()
	sum := 0.0
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
