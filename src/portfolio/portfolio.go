package portfolio

import (
	"errors"
	"fmt"
	"github.com/stock-simulator-server/src/valuable"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/duplicator"
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
	PersonalLedger map[string]*ledgerEntry `json:"ledger" change:"-"`

	UpdateChannel   *duplicator.ChannelDuplicator `json:"-"`
	valuableUpdates *duplicator.ChannelDuplicator `json:"-"`

	Lock *lock.Lock `json:"-"`
}

type ledgerEntry struct {
	Amount        float64 `json:"amount"`
	updateChannel chan interface{} `json:"-"`
}

func (port *Portfolio) GetId() string {
	return port.UUID
}

func (port *Portfolio) GetType() string {
	return ObjectType
}

func NewPortfolio(userUUID, name string) (*Portfolio, error) {
	PortfoliosLock.Acquire("new-portfolio")
	defer PortfoliosLock.Release()
	if _, exists := Portfolios[userUUID]; exists {
		return nil, errors.New("portfolio uuid already Exists")
	}
	port :=
		&Portfolio{
			Name:            name,
			UUID:            userUUID,
			Wallet:          1000,
			UpdateChannel:   duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-update", userUUID)),
			Lock:            lock.NewLock(fmt.Sprintf("portfolio-%s", name)),
			valuableUpdates: duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-valueable-update", userUUID)),
			PersonalLedger:  make(map[string]*ledgerEntry),
		}
	Portfolios[userUUID] = port
	PortfoliosUpdateChannel.RegisterInput(port.UpdateChannel.GetOutput())
	go port.valuableUpdate()
	NewPortfolioChannel.Offer(port)
	return port, nil
}
func (port *Portfolio) valuableUpdate() {
	updateChannel := port.valuableUpdates.GetOutput()

	for range updateChannel {
		port.Lock.Acquire("portfolio-update")
		newNetWorth := port.calculateNetWorth()
		if newNetWorth != port.NetWorth {
			port.NetWorth = newNetWorth
			port.UpdateChannel.Offer(port)
		}
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
	sum := 0.0
	for valueStr, entry := range port.PersonalLedger {
		value := valuable.Stocks[valueStr]
		sum += value.GetValue() * entry.Amount
	}
	return sum + port.Wallet
}

func (port *Portfolio) TradeUpdate(value valuable.Valuable, amountOwned, price float64) {
	valueID := value.GetId()
	entry, exists := port.PersonalLedger[valueID]
	if !exists {
		port.PersonalLedger[valueID] = &ledgerEntry{
			Amount:        amountOwned,
			updateChannel: value.GetUpdateChannel().GetOutput(),
		}
		entry = port.PersonalLedger[valueID]
		port.valuableUpdates.RegisterInput(entry.updateChannel)
	}

	if amountOwned == 0 {
		value.GetUpdateChannel().UnregisterOutput(entry.updateChannel)
		close(entry.updateChannel)
		delete(port.PersonalLedger, valueID)
	} else {
		entry.Amount = amountOwned
	}

	port.Wallet -= price
	port.NetWorth = port.calculateNetWorth()
	port.UpdateChannel.Offer(port)

}

func GetAllPortfolios()[]*Portfolio{
	PortfoliosLock.Acquire("get all ports")
	defer PortfoliosLock.Release()
	lst := make([]*Portfolio, len(Portfolios))
	i := 0
	for _, val := range Portfolios{
		lst[i] = val
		i+= 1
	}
	return lst
}