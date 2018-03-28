package valuable

import (
	"errors"
	"fmt"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/utils"
	"math/rand"
	"reflect"
	"time"
)

const (
	volatilityMin      = 1
	volatilityMax      = 10
	volatilityMinTurns = 1
	volatilityMaxTurns = 100

	timeSimulationPeriod = time.Second

	ObjectType = "stock"
)

var timeSimulation = duplicator.MakeDuplicator("time-sim")

var Stocks = make(map[string]*Stock)
var NewStockChannel = duplicator.MakeDuplicator("new-stock-channel")

func StartStockStimulation() {
	ticker := time.NewTicker(timeSimulationPeriod)
	simulation := make(chan interface{})
	timeSimulation.RegisterInput(simulation)
	go func() {
		for range ticker.C {
			simulation <- true
		}
		close(simulation)
	}()

}

type stockManager struct {
	stocks             map[string]*Stock
	StockUpdateChannel *duplicator.ChannelDuplicator
}

//Stock type for storing the stock information
type Stock struct {
	Uuid          string                        `json:"uuid"`
	Name          string                        `json:"name"`
	TickerId      string                        `json:"ticker_id"`
	CurrentPrice  float64                       `json:"current_price" change:"-"`
	OpenShares    float64                       `json:"open_shares"`
	PriceChanger  PriceChange                   `json:"-"`
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
	lock          *lock.Lock                    `json:"-"`
}

func (stock *Stock) GetType() string {
	return ObjectType
}

func NewStock(tickerID, name string, startPrice float64, runInterval time.Duration) (*Stock, error) {
	// Acquire the valuableMapLock so no one can add a new entry till we are done
	ValuablesLock.Acquire("new-stock")
	defer ValuablesLock.Release()
	for _, s := range Stocks {
		if s.TickerId == tickerID {
			return nil, errors.New("tickerID is already taken by another valuable")
		}
	}
	uuidString := utils.PseudoUuid()
	for {
		if _, ok := Stocks[uuidString]; !ok {
			break
		}
		uuidString = utils.PseudoUuid()
	}

	stock := &Stock{
		Uuid:          uuidString,
		lock:          lock.NewLock(fmt.Sprintf("stock-%s", tickerID)),
		Name:          name,
		TickerId:      tickerID,
		CurrentPrice:  startPrice,
		UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("stock-%s-update", tickerID)),
	}

	stock.PriceChanger = &RandomPrice{
		RunPercent:            timeSimulationPeriod.Seconds() / (runInterval.Seconds() * 1.0),
		TargetPrice:           100.0,
		PercentToChangeTarget: .1,
		Volatility:            5,
	}
	go stock.stockUpdateRoutine()
	Stocks[tickerID] = stock
	ValuableUpdateChannel.RegisterInput(stock.UpdateChannel.GetOutput())
	ValuableUpdateChannel.Offer(stock)
	//NewStockChannel.Offer(stock)
	return stock, nil
}

func (stock *Stock) GetValue() float64 {
	return stock.CurrentPrice
}

func (stock *Stock) GetName() string {
	return stock.Name
}

func (stock *Stock) GetLock() *lock.Lock {
	return stock.lock
}

func (stock *Stock) GetUpdateChannel() *duplicator.ChannelDuplicator {
	return stock.UpdateChannel
}

func (stock *Stock) GetId() string {
	return stock.TickerId
}

func (stock *Stock) stockUpdateRoutine() {
	update := timeSimulation.GetOutput()
	for range update {
		stock.PriceChanger.change(stock)
	}
}

func (stock *Stock) ChangeDetected() reflect.Type {
	return reflect.TypeOf(stock)
}

// Some thing that can take in a stock and change the current price
type PriceChange interface {
	change(stock *Stock)
}

// Random Price implements priceChange
type RandomPrice struct {
	RunPercent            float64 `json:"run_percent"`
	TargetPrice           float64 `json:"target_price"`
	PercentToChangeTarget float64 `json:"change_percent"`
	Volatility            float64 `json:"volatility"`
}

//change the stock using the changer
func (randPrice *RandomPrice) change(stock *Stock) {
	stock.lock.Acquire("change-stock")
	defer stock.lock.Release()

	if rand.Float64() >= randPrice.RunPercent {
		return
	}
	if rand.Float64() <= randPrice.PercentToChangeTarget {
		randPrice.changeValues()
	}

	//can make this a lot more interesting, like adding in the ability for it to drop
	change := (randPrice.TargetPrice - stock.CurrentPrice) /
		utils.MapNum(randPrice.Volatility, volatilityMin, volatilityMax, volatilityMinTurns, volatilityMaxTurns)

	stock.CurrentPrice = stock.CurrentPrice + change

	stock.UpdateChannel.Offer(stock)

}

// change the price of the changer
func (randPrice *RandomPrice) changeValues() {

	// get what the upper and lower bounds in % of the current price
	window := utils.MapNum(randPrice.Volatility, volatilityMin, volatilityMax, 0, 0.3)
	// select a random number on +- that %
	newTarget := utils.MapNum(rand.Float64(), 0, 1, randPrice.TargetPrice*(1-window), randPrice.TargetPrice*(1+window))
	//need to deiced if the floor should happen before or after
	if newTarget < 0 {
		newTarget = 0
	}

	randPrice.TargetPrice = newTarget
	randPrice.Volatility = utils.RandRange(volatilityMin, volatilityMax)
}

func GetAllStocks() []*Stock {
	ValuablesLock.Acquire("Get List")
	defer ValuablesLock.Release()
	v := make([]*Stock, len(Stocks))
	i := 0
	for _, val := range Stocks {
		v[i] = val
		i += 1
	}
	return v

}

/** ########################################
*           Math Helper Functions
*   ########################################
 */
