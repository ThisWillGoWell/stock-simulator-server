package valuable

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/wires"
)

const (
	volatilityMin      = 1
	volatilityMax      = 10
	volatilityMinTurns = 5
	volatilityMaxTurns = 25

	timeSimulationPeriod = time.Second

	ObjectType = "stock"
)

var timeSimulation = duplicator.MakeDuplicator("time-sim")

var Stocks = make(map[string]*Stock)

func StartStockStimulation() {
	/*
		when simulation gets emitted, it will trigger all the tings listening to it
		to run a session
	*/
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

//Stock type for storing the stock information
type Stock struct {
	Uuid           string                        `json:"uuid"`
	Name           string                        `json:"name"`
	TickerId       string                        `json:"ticker_id"`
	CurrentPrice   int64                         `json:"current_price" change:"-"`
	OpenShares     int64                         `json:"open_shares" change:"-"`
	ChangeDuration time.Duration                 `json:"-"`
	PriceChanger   PriceChange                   `json:"-"`
	UpdateChannel  *duplicator.ChannelDuplicator `json:"-"`
	lock           *lock.Lock                    `json:"-"`
}

func (stock *Stock) GetType() string {
	return ObjectType
}

func NewStock(tickerID, name string, startPrice int64, runInterval time.Duration) (*Stock, error) {
	// Acquire the valuableMapLock so no one can add a new entry till we are done
	ValuablesLock.Acquire("new-stock")
	defer ValuablesLock.Release()

	uuidString := utils.SerialUuid()

	return MakeStock(uuidString, tickerID, name, startPrice, 100, runInterval)
}

func MakeStock(uuid, tickerID, name string, startPrice, openShares int64, runInterval time.Duration) (*Stock, error) {
	for _, s := range Stocks {
		if s.TickerId == tickerID {
			return nil, errors.New("tickerID is already taken by another valuable")
		}
	}
	stock := &Stock{
		ChangeDuration: runInterval,
		OpenShares:     openShares,
		Uuid:           uuid,
		lock:           lock.NewLock(fmt.Sprintf("stock-%s", tickerID)),
		Name:           name,
		TickerId:       tickerID,
		CurrentPrice:   startPrice,
		UpdateChannel:  duplicator.MakeDuplicator(fmt.Sprintf("stock-%s-update", tickerID)),
	}
	//stock.lock.EnableDebug()

	stock.PriceChanger = &RandomPrice{
		RunPercent:            timeSimulationPeriod.Seconds() / (runInterval.Seconds() * 1.0),
		TargetPrice:           int64(rand.Intn(100000)),
		PercentToChangeTarget: .07,
		Volatility:            5,
		RandomNoise:           .07,
	}
	go stock.stockUpdateRoutine()
	Stocks[uuid] = stock
	stock.UpdateChannel.EnableCopyMode()
	change.RegisterPublicChangeDetect(stock)
	wires.StocksUpdate.RegisterInput(stock.UpdateChannel.GetBufferedOutput(1000))
	wires.StocksNewObject.Offer(stock)
	utils.RegisterUuid(uuid, stock)
	return stock, nil
}

func (stock *Stock) GetValue() int64 {
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
	return stock.Uuid
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
	GetTargetPrice() int64
}

// Random Price implements priceChange
type RandomPrice struct {
	RunPercent            float64 `json:"run_percent"`
	TargetPrice           int64   `json:"target_price"`
	PercentToChangeTarget float64 `json:"change_percent"`
	Volatility            float64 `json:"volatility"`
	RandomNoise           float64
}

func (randPrice *RandomPrice) GetTargetPrice() int64 {
	return randPrice.TargetPrice
}

//change the stock using the changer
func (randPrice *RandomPrice) change(stock *Stock) {
	if rand.Float64() >= randPrice.RunPercent {
		return
	}
	stock.lock.Acquire("change-stock")
	defer stock.lock.Release()

	if rand.Float64() <= randPrice.PercentToChangeTarget {
		randPrice.changeValues()
	}

	moveToTarget := int64(utils.RandRangeFloat(float64(randPrice.TargetPrice)*0.9, float64(randPrice.TargetPrice)*1.1))

	change := float64(moveToTarget-stock.CurrentPrice) /
		utils.MapNumFloat(randPrice.Volatility, volatilityMin, volatilityMax, volatilityMinTurns, volatilityMaxTurns)

	if rand.Float64() <= randPrice.RandomNoise {
		change = change * -1
	}
	stock.CurrentPrice = int64(float64(stock.CurrentPrice) + (change * .5))

	stock.UpdateChannel.Offer(stock)

}

// change the price of the changer
func (randPrice *RandomPrice) changeValues() {

	// get what the upper and lower bounds in % of the current price
	window := utils.MapNumFloat(randPrice.Volatility, volatilityMin, volatilityMax, 0, 0.3)
	// select a random number on +- that
	newTarget := utils.MapNumFloat(rand.Float64(), 0, 1, float64(randPrice.TargetPrice)*(1-window), float64(randPrice.TargetPrice)*(1+window))
	// this is to prevent all the stocks to becoming 1.231234e-14
	if newTarget < 500 {
		if rand.Float64() < randPrice.PercentToChangeTarget {
			newTarget = 1000 + newTarget
		}
	}
	//need to deiced if the floor should happen before or after
	if newTarget < 0 {
		newTarget = 0
	}

	randPrice.TargetPrice = int64(newTarget)
	randPrice.Volatility = utils.RandRangeFloat(volatilityMin, volatilityMax)
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

func (s *Stock) Update() {
	s.UpdateChannel.Offer(s)
}

/** ########################################
*           Math Helper Functions
*   ########################################
 */
