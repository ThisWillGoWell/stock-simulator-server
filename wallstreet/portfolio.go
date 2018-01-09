package wallstreet

import "stock-server/utils"

var PortfolioIds map[string]*Portfolio


type Portfolio struct {
	Name     string
	UUID     string
	Exchange map[string]Exchange
	Wallet   float64
	NetWorth float64
	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	personalLedger map[*Stock]float64


	UpdateChannel *utils.ChannelDuplicator
	stockUpdates  *utils.ChannelDuplicator

	lock *utils.Lock
}

func NewPortfolio( userUUID,  name string )(*Portfolio){
	return &Portfolio{
		Name:           name,
		UUID:           userUUID,
		Wallet:         1000,
		UpdateChannel:  utils.MakeDuplicator(),
		lock :          utils.NewLock(),
	}
}

func (port *Portfolio) updateWallet(walletChange float64){
	port.NetWorth += walletChange
}