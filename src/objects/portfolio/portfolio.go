package portfolio

import (
	"errors"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/ThisWillGoWell/stock-simulator-server/src/id"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/game/money"

	"github.com/ThisWillGoWell/stock-simulator-server/src/game/level"
	"github.com/ThisWillGoWell/stock-simulator-server/src/id/change"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/ledger"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/duplicator"
)

var Portfolios = make(map[string]*Portfolio)
var PortfoliosLock = lock.NewLock("portfolios")

/**
Portfolios are the $$$ part of a user
*/
type Portfolio struct {
	objects.Portfolio
	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	// stock_uuid -> ledgerObject
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
	UpdateInput   *duplicator.ChannelDuplicator `json:"-"`
	Lock          *lock.Lock                    `json:"-"`
	close         chan interface{}
}

func NewPortfolio(portfolioUuid, userUuid string) (*Portfolio, error) {
	PortfoliosLock.Acquire("new-portfolio")
	defer PortfoliosLock.Release()

	portfolio := objects.Portfolio{
		UserUUID: userUuid,
		Uuid:     portfolioUuid,
		Wallet:   10 * money.Thousand,
		NetWorth: 10 * money.Thousand,
		Level:    0,
	}
	port, err := MakePortfolio(portfolio, true)
	if err != nil {
		return nil, err
	}
	return port, err
}

func DeletePortfolio(uuid string) {
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

func MakePortfolio(portfolio objects.Portfolio, lockAquired bool) (*Portfolio, error) {
	//PortfoliosUpdateChannel.EnableDebug("port update")
	if !lockAquired {
		PortfoliosLock.Acquire("new-portfolio")
		defer PortfoliosLock.Release()
	}
	if _, exists := Portfolios[portfolio.Uuid]; exists {
		return nil, errors.New("portfolio uuid already Exists")
	}
	port :=
		&Portfolio{
			Portfolio:     portfolio,
			UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-update", portfolio.Uuid)),
			Lock:          lock.NewLock(fmt.Sprintf("portfolio-%s", portfolio.Uuid)),
			UpdateInput:   duplicator.MakeDuplicator(fmt.Sprintf("portfolio-%s-valueable-update", portfolio.Uuid)),
		}
	port.Lock.EnableDebug()
	port.UpdateChannel.EnableCopyMode()
	if err := change.RegisterPublicChangeDetect(port); err != nil {
		return nil, err
	}
	Portfolios[port.Uuid] = port
	wires.PortfolioUpdate.RegisterInput(port.UpdateChannel.GetBufferedOutput(1000))
	id.RegisterUuid(port.Uuid, port)
	go port.valuableUpdate()
	return port, nil
}

/**
async code that gets called whenever a stock or a ledger that the portfolio owns changes
This then triggers a recalc of net worth and offers its self up as a update
*/
func (port *Portfolio) valuableUpdate() {
	updateChannel := port.UpdateInput.GetBufferedOutput(100)
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
	port.UpdateChannel.Offer(port.Portfolio)
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
	effect.EffectLock.Acquire("level-up")
	defer effect.EffectLock.Release()
	newEffect, oldEffect, err := effect.UpdateBaseProfit(port.Uuid, l.ProfitMultiplier)
	if err != nil {
		log.Log.Errorf("failed to level up, err new effect err=[%v]", err)
		return fmt.Errorf("opps! something went wrong 0x24")
	}

	port.Wallet = port.Wallet - l.Cost
	port.Level = nextLevel

	// commit to database
	if dbErr := database.Db.Execute([]interface{}{newEffect, port}, []interface{}{oldEffect}); dbErr != nil {
		log.Log.Errorf("failed to write level up to database err=[%v]", err)
		return fmt.Errorf("opps! something went wrong! 0x34")
	}

	wires.EffectsDelete.Offer(newEffect.Effect)
	effect.DeleteEffect(oldEffect)
	wires.EffectsDelete.Offer(oldEffect.Effect)

	go port.Update()
	return nil
}
