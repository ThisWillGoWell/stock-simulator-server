package valuable

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

	"github.com/ThisWillGoWell/stock-simulator-server/src/change"

	"github.com/ThisWillGoWell/stock-simulator-server/src/money"

	"github.com/ThisWillGoWell/stock-simulator-server/src/duplicator"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
)

const (
	volatilityMin      = 1
	volatilityMax      = 10
	volatilityMinTurns = 10
	volatilityMaxTurns = 100

	timeSimulationPeriod = time.Millisecond * 500

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
	models.Stock
	PriceChanger  PriceChange                   `json:"-"`
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
	lock          *lock.Lock                    `json:"-"`
	close         chan interface{}              `json:"-"`
}

func (stock *Stock) GetType() string {
	return ObjectType
}

func NewStock(tickerID, name string, startPrice int64, runInterval time.Duration) (*Stock, error) {
	// Acquire the valuableMapLock so no one can add a new entry till we are done
	ValuablesLock.Acquire("new-stock")
	defer ValuablesLock.Release()
	stock := models.Stock{
		OpenShares:     1000,
		Uuid:           id.SerialUuid(),
		Name:           name,
		TickerId:       tickerID,
		CurrentPrice:   startPrice,
		ChangeDuration: runInterval,
	}
	s, err := MakeStock(stock)
	if err != nil {
		return nil, err
	}
	wires.StocksNewObject.Offer(s)
	return s, nil
}

func DeleteStock(s *Stock) {
	s.UpdateChannel.StopDuplicator()
	close(s.close)
	delete(Stocks, s.Uuid)
	id.RemoveUuid(s.Uuid)
	change.UnregisterChangeDetect(s)
}

func MakeStock(s models.Stock) (*Stock, error) {
	for _, currentStocks := range Stocks {
		if currentStocks.TickerId == s.TickerId {
			return nil, errors.New("tickerID is already taken by another valuable")
		}
	}

	stock := &Stock{
		Stock: s,
		lock:          lock.NewLock(fmt.Sprintf("stock-%s", s.TickerId)),
		UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("stock-%s-update", s.TickerId)),
	}
	//stock.lock.EnableDebug()

	stock.PriceChanger = &RandomPrice{
		RunPercent:            timeSimulationPeriod.Seconds() / (s.ChangeDuration.Seconds() * 1.0),
		TargetPrice:           int64(rand.Intn(100000)),
		PercentToChangeTarget: .1,
		Volatility:            5,
		RandomNoise:           .15,
	}
	if err := change.RegisterPublicChangeDetect(stock); err != nil {
		return nil, err
	}
	go stock.stockUpdateRoutine()
	Stocks[s.Uuid] = stock
	stock.UpdateChannel.EnableCopyMode()
	wires.StocksUpdate.RegisterInput(stock.UpdateChannel.GetBufferedOutput(1000))
	id.RegisterUuid(s.Uuid, stock)
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
loop:
	for {
		select {
		case <-update:
			stock.PriceChanger.change(stock)
		case <-stock.close:
			break loop
		}
	}
	timeSimulation.UnregisterOutput(update)
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
		change = change * -1 * .25
	}
	startPrice := stock.CurrentPrice
	stock.CurrentPrice = int64(float64(stock.CurrentPrice) + (change * .5))
	if err := database.Db.Execute([]interface{}{stock}, nil); err != nil {
		log.Log.Errorf("failed to update stock err=[%v]", err)
		stock.CurrentPrice = startPrice
		return
	}

	stock.UpdateChannel.Offer(stock)

}

// change the price of the changer
func (randPrice *RandomPrice) changeValues() {

	// get what the upper and lower bounds in % of the current price
	window := utils.MapNumFloat(randPrice.Volatility, volatilityMin, volatilityMax, 0, 0.4)
	// select a random number on +- that
	newTarget := utils.MapNumFloat(rand.Float64(), 0, 1, float64(randPrice.TargetPrice)*(1-window), float64(randPrice.TargetPrice)*(1+window))
	// this is to prevent all the stocks to becoming 1.231234e-14
	if newTarget < 500 {
		if rand.Float64() < randPrice.PercentToChangeTarget {
			newTarget = 1000 + newTarget
		}
	}
	// to add
	if newTarget < 100*money.Thousand && rand.Float64() < .05 {
		if rand.Float64() < .5 {
			newTarget = newTarget * .5
		} else {
			newTarget = newTarget * 2
		}
	}
	if newTarget > 1*money.Million {
		newTarget = 1*money.Million + (newTarget-1*money.Million)/2
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
