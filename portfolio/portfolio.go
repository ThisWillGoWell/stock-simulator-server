package portfolio

import (
	"stock-server/utils"
	"stock-server/valuable"
	"errors"
	"fmt"
)

var Portfolios = make(map[string]*Portfolio)
var PortfoliosLock = utils.NewLock()


type Portfolio struct {
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Wallet   float64 `json:"wallet"`
	NetWorth float64 `json:"net_worth"`

	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	PersonalLedger map[valuable.Valuable]*ledgerEntry `json:"ledger"`

	UpdateChannel   *utils.ChannelDuplicator
	valuableUpdates *utils.ChannelDuplicator

	Lock *utils.Lock
}

type ledgerEntry struct {
	amount float64
	updateChannel chan interface{}
}

func NewPortfolio( userUUID string)(*Portfolio, error){
	PortfoliosLock.Acquire()
	defer PortfoliosLock.Release()
	if _, exists := Portfolios[userUUID]; exists {
		return nil, errors.New("portfolio uuid already Exists")
	}
	 port :=
		&Portfolio{
			UUID:           userUUID,
			Wallet:         1000,
			UpdateChannel:  utils.MakeDuplicator(),
			Lock :          utils.NewLock(),
			valuableUpdates: utils.MakeDuplicator(),
		}
	Portfolios[userUUID] = port
	go port.valuableUpdate()
	return port, nil
}
func (port *Portfolio)valuableUpdate(){
	updateChannel := port.valuableUpdates.GetOutput()

	for range updateChannel{
		fmt.Println("port got update")
		port.Lock.Acquire()
		port.NetWorth = port.calculateNetWorth()
		port.Lock.Release()
		port.UpdateChannel.Offer(port)
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
	for value, entry := range port.PersonalLedger{
		sum += value.GetValue() * entry.amount
	}
	return sum + port.Wallet
}

func (port *Portfolio)UpdateLedger(value valuable.Valuable, amount float64){
	entry, exists := port.PersonalLedger[value]
	if !exists{
		port.PersonalLedger[value] = &ledgerEntry{
			amount: amount,
			updateChannel: value.GetUpdateChannel().GetOutput(),
		}
		entry = port.PersonalLedger[value]
		port.valuableUpdates.RegisterInput(entry.updateChannel)
	}

	if amount == 0{
		value.GetUpdateChannel().UnregisterOutput(entry.updateChannel)
		close(entry.updateChannel)
		delete(port.PersonalLedger, value)
	} else{
		entry.amount = amount
	}
	port.NetWorth = port.calculateNetWorth()

}


func (port *Portfolio) UpdateWallet(walletChange float64){
	port.NetWorth += walletChange
}