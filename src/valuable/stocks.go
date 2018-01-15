package valuable

import (
	"math"
	"math/rand"
	"fmt"
	"github.com/stock-simulator-server/src/utils"
	"time"
	"errors"
)

const (
	volatilityMin      = 1
	volatilityMax      = 10
	volatilityMinTurns = 1
	volatilityMaxTurns = 100

	timeSimulationPeriod = time.Second


)
var timeSimulation = utils.MakeDuplicator()

var manager = stockManager{
	stocks: make(map[string]*Stock),
	StockUpdateChannel: utils.MakeDuplicator(),
}

func StartStockStimulation(){
	ticker := time.NewTicker(timeSimulationPeriod)
	simulation := make(chan interface{})
	timeSimulation.RegisterInput(simulation)
	go func(){
		for range ticker.C {
			simulation <- true
		}
		close(simulation)
	}();

}

type stockManager struct {
	stocks             map[string]*Stock
	StockUpdateChannel *utils.ChannelDuplicator
}
//Stock type for storing the stock information
type Stock struct {
	Name         string  `json:"name"`
	TickerId     string  `json:"ticker_id"`
	CurrentPrice float64 `json:"current_price"`
	PriceChanger PriceChange `json:"-"`
	UpdateChannel *utils.ChannelDuplicator `json:"-"`
	lock *utils.Lock
}

func NewStock(tickerID, name string, startPrice float64, runInterval time.Duration)(*Stock, error){
	// Acquire the valuableMapLock so no one can add a new entry till we are done
	ValuablesLock.Acquire("new-stock")
	defer ValuablesLock.Release()
	if _, ok := Valuables[tickerID]; ok{
		return nil, errors.New("tickerID is already taken by another valuable")
	}
	stock := &Stock{
		lock: utils.NewLock(fmt.Sprintf("stock-%s", tickerID)),
		Name: name,
		TickerId: tickerID,
		CurrentPrice: startPrice,
		UpdateChannel: utils.MakeDuplicator(),
	}

	stock.PriceChanger = &RandomPrice{
		RunPercent:            timeSimulationPeriod.Seconds() / (runInterval.Seconds() * 1.0),
		TargetPrice:           100.0,
		PercentToChangeTarget: .1,
		Volatility:            5,
	}
	go stock.stockUpdateRoutine()
	Valuables[tickerID] = stock
	ValuableUpdateChannel.RegisterInput(stock.UpdateChannel.GetOutput())
	return stock, nil
}

func (stock *Stock)GetValue()float64{
	return stock.CurrentPrice
}

func (stock *Stock)GetLock() *utils.Lock{
	return stock.lock
}

func  (stock *Stock)GetUpdateChannel()(*utils.ChannelDuplicator){
	return stock.UpdateChannel
}

func (stock *Stock)GetID() string{
	return stock.TickerId
}

func (stock *Stock)stockUpdateRoutine(){
	update := timeSimulation.GetOutput()
	for range update{
		stock.PriceChanger.change(stock)
	}
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
func (randPrice *RandomPrice) change(stock *Stock){
	stock.lock.Acquire("change-stock")
	defer stock.lock.Release()

	if rand.Float64() >= randPrice.RunPercent {
		return
	}
	if rand.Float64() <= randPrice.PercentToChangeTarget {
		randPrice.changeValues()
	}

	//can make this a lot more interesting, like adding in the ability for it to drop
	change :=  (randPrice.TargetPrice - stock.CurrentPrice) /
		MapNum(randPrice.Volatility, volatilityMin, volatilityMax, volatilityMinTurns, volatilityMaxTurns)

	stock.CurrentPrice = stock.CurrentPrice + change

	stock.UpdateChannel.Offer(stock)

}

// change the price of the changer
func (randPrice *RandomPrice)changeValues(){

	// get what the upper and lower bounds in % of the current price
	window := MapNum(randPrice.Volatility, volatilityMin, volatilityMax, 0, 0.3)
	// select a random number on +- that %
	newTarget := MapNum(rand.Float64(), 0, 1, randPrice.TargetPrice * (1 - window), randPrice.TargetPrice * (1 + window))
	fmt.Println("old", randPrice.TargetPrice, "new", newTarget, "window", window)
	//need to deiced if the floor should happen before or after
	if newTarget < 0{
		newTarget = 0
	}

	randPrice.TargetPrice = newTarget
	randPrice.Volatility = RandRange(volatilityMin, volatilityMax)
}

/** ########################################
*           Math Helper Functions
*   ########################################
*/

// Round f to nearest number of decimal points
func RoundPlus(f float64, places int) (float64) {
	shift := math.Pow(10, float64(places))
	return Round(f * shift) / shift
}

// Round a float to the nearest int
func Round(f float64) float64 {
	return math.Floor(f + .5)
}

// generate a random number between two floats
func RandRange(min, max float64) float64 {
	return MapNum(rand.Float64(), 0, 1, min, max)
}


// map a number from one range to another range
func MapNum(value, inMin, inMax, outMin, outMax float64) float64{
	if value >= inMax{
		return outMax
	}
	if value <= inMin{
		return outMin
	}
	return (value - inMin) * (outMax - outMin) / (inMax - inMin) + outMin
}

