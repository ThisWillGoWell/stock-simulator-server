package portfolio

import (
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
	"errors"
	"fmt"
)

const(
	ObjectType = "Portfolio"
)

var Portfolios = make(map[string]*Portfolio)
var PortfoliosLock = utils.NewLock("portfolios")
var PortfoliosUpdateChannel = utils.MakeDuplicator()

type Portfolio struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Wallet   float64 `json:"wallet" change:"-"`
	NetWorth float64 `json:"net_worth" change:"-"`

	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	PersonalLedger map[string]*ledgerEntry `json:"ledger" change:"-"`

	UpdateChannel   *utils.ChannelDuplicator `json:"-"`
	valuableUpdates *utils.ChannelDuplicator

	Lock *utils.Lock `json:"-"`
}

type ledgerEntry struct {
	Amount float64 `json:"amount"`
	updateChannel chan interface{}
}

func(port *Portfolio)GetId() string{
	return port.UUID
}


func(port *Portfolio)GetType() string{
	return ObjectType
}

func NewPortfolio( userUUID, name string)(*Portfolio, error){
	PortfoliosLock.Acquire("new-portfolio")
	defer PortfoliosLock.Release()
	if _, exists := Portfolios[userUUID]; exists {
		return nil, errors.New("portfolio uuid already Exists")
	}
	 port :=
		&Portfolio{
			Name:			name,
			UUID:           userUUID,
			Wallet:         1000,
			UpdateChannel:  utils.MakeDuplicator(),
			Lock :          utils.NewLock(fmt.Sprintf("portfolio-%s", name)),
			valuableUpdates: utils.MakeDuplicator(),
			PersonalLedger: make(map[string]*ledgerEntry),
		}
	Portfolios[userUUID] = port
	PortfoliosUpdateChannel.RegisterInput(port.UpdateChannel.GetOutput())
	go port.valuableUpdate()
	return port, nil
}
func (port *Portfolio)valuableUpdate(){
	updateChannel := port.valuableUpdates.GetOutput()

	for range updateChannel{
		port.Lock.Acquire("portfolio-update")
		newNetWorth := port.calculateNetWorth()
		if newNetWorth != port.NetWorth{
			port.NetWorth = newNetWorth
			port.UpdateChannel.Offer(port)
		}
		port.Lock.Release()
	}
}

func GetPortfolio(userUUID string)(*Portfolio, error){
	port, exists := Portfolios[userUUID]
	if !exists {
		return nil, errors.New("uuid does not have a portfolio tied to it")
	}
	return port, nil
}


//update the current net worth. NOT THREAD SAFE
func (port *Portfolio)calculateNetWorth() float64{
	sum := 0.0
	for valueStr, entry := range port.PersonalLedger{
		value :=valuable.Valuables[valueStr]
		sum += value.GetValue() * entry.Amount
	}
	return sum + port.Wallet
}

func (port *Portfolio)TradeUpdate(value valuable.Valuable, amountOwned, price float64){
	valueID := value.GetId()
	entry, exists := port.PersonalLedger[valueID]
	if !exists{
		port.PersonalLedger[valueID] = &ledgerEntry{
			Amount: amountOwned,
			updateChannel: value.GetUpdateChannel().GetOutput(),
		}
		entry = port.PersonalLedger[valueID]
		port.valuableUpdates.RegisterInput(entry.updateChannel)
	}

	if amountOwned == 0{
		value.GetUpdateChannel().UnregisterOutput(entry.updateChannel)
		close(entry.updateChannel)
		delete(port.PersonalLedger, valueID)
	} else{
		entry.Amount = amountOwned
	}

	port.Wallet -= price
	port.NetWorth = port.calculateNetWorth()
	port.UpdateChannel.Offer(port)

}