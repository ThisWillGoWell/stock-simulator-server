package portfolio

import (
	"errors"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

	"github.com/ThisWillGoWell/stock-simulator-server/src/effect"

	"github.com/ThisWillGoWell/stock-simulator-server/src/money"

	"github.com/ThisWillGoWell/stock-simulator-server/src/change"
	"github.com/ThisWillGoWell/stock-simulator-server/src/duplicator"
	"github.com/ThisWillGoWell/stock-simulator-server/src/ledger"
	"github.com/ThisWillGoWell/stock-simulator-server/src/level"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/valuable"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
)

const (
	ObjectType = "portfolio"
)

var Portfolios = make(map[string]*Portfolio)
var PortfoliosLock = lock.NewLock("portfolios")

/**
Portfolios are the $$$ part of a user
*/
type Portfolio struct {
	models.Portfolio
	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	// stock_uuid -> ledgerObject
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
	UpdateInput   *duplicator.ChannelDuplicator `json:"-"`
	Lock          *lock.Lock                    `json:"-"`
	close         chan interface{}
}

func (port *Portfolio) GetId() string {
	return port.Uuid
}

func (port *Portfolio) GetType() string {
	return ObjectType
}

func NewPortfolio(portfolioUuid, userUuid string) (*Portfolio, error) {
	PortfoliosLock.Acquire("new-portfolio")
	defer PortfoliosLock.Release()

	port, err := MakePortfolio(portfolioUuid, userUuid, 10*money.Thousand, 0, true)
	if err != nil {
		return nil, err
	}
	return port, err
}

func DeletePortfolio(uuid string, lockAquired, force bool) {
	if !lockAquired {
		PortfoliosLock.Acquire("delete-portfolio")
		defer PortfoliosLock.Release()
	}
	port, ok := Portfolios[uuid]
	if !ok {
		log.Log.Errorf("ot a portfolio delete on a uuid not found")
		return
	}
	close(port.close)
	port.UpdateInput.StopDuplicator()
	port.UpdateChannel.StopDuplicator()
	change.UnregisterChangeDetect(port)
	delete(Portfolios, uuid)
	id.RemoveUuid(uuid)
}

func MakePortfolio(uuid, userUUID string, wallet, level int64, lockAquired bool) (*Portfolio, error) {
	//PortfoliosUpdateChannel.EnableDebug("port update")
	if !lockAquired {
		PortfoliosLock.Acquire("new-portfolio")
		defer PortfoliosLock.Release()
	}
	if _, exists := Portfolios[uuid]; exists {
		id.RemoveUuid(uuid)
		return nil, errors.New("portfolio uuid already Exists")
	}
	port :=
		&Portfolio{
			Portfolio: models.Portfolio{
				UserUUID: userUUID,
				Uuid:     uuid,
				Wallet:   wallet,
				NetWorth: wallet,
				Level:    level,
			},
			UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-update", uuid)),
			Lock:          lock.NewLock(fmt.Sprintf("portfolio-%s", uuid)),
			UpdateInput:   duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-valueable-update", uuid)),
		}

	port.UpdateChannel.EnableCopyMode()
	if err := change.RegisterPublicChangeDetect(port); err != nil {
		return nil, err
	}
	Portfolios[uuid] = port

	wires.PortfolioNewObject.Offer(port)
	wires.PortfolioUpdate.RegisterInput(port.UpdateChannel.GetBufferedOutput(1000))
	id.RegisterUuid(uuid, port)
	go port.valuableUpdate()
	return port, nil
}

/**
async code that gets called whenever a stock or a ledger that the portfolio owns changes
This then triggers a recalc of net worth and offers its self up as a update
*/
func (port *Portfolio) valuableUpdate() {
	updateChannel := port.UpdateInput.GetBufferedOutput(1000)
	for range updateChannel {
		port.Update()
	}
}

func UpdateAll() {
	for _, p := range Portfolios {
		p.Update()
	}
}

func (port *Portfolio) Update() {
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
	for valueStr, entry := range ledger.EntriesPortfolioStock[port.Uuid] {
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

func (port *Portfolio) LevelUp() error {
	port.Lock.Acquire("level up")
	defer port.Lock.Release()
	nextLevel := port.Level + 1
	l, exists := level.Levels[nextLevel]
	if !exists {
		return errors.New("there is no next l")
	}
	if port.Wallet < l.Cost {
		return errors.New("not enough $$")
	}
	port.Wallet = port.Wallet - l.Cost
	port.Level = nextLevel
	effect.UpdateBaseProfit(port.Uuid, l.ProfitMultiplier)
	go port.Update()
	return nil
}
