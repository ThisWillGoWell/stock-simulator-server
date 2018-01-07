package wallstreet

import (
	"math"
	"math/rand"
	"fmt"
	"stock-server/utils"
)

const (
	volatilityMin      = 1
	volatilityMax      = 10
	volatilityMinTurns = 1
	volatilityMaxTurns = 100

)

type StockManager struct {
	stocks             map[string]*Stock
	StockUpdateChannel utils.ChannelDuplicator
}


func (manager *StockManager) changeStock(){
	for _ , stock := range manager.stocks {
		stock.PriceChanger.change(stock)
	}
}

func (manager *StockManager) getStock(ticker string)(*Stock){
	return manager.stocks[ticker]
}

func (manager *StockManager) addStock(tickerId, name string, startPrice , runPercent float64) *Stock{
	if _, ok := manager.stocks[tickerId]; ok{
		return nil
	}
	stock := Stock{
		Name: name,
		TickerId: tickerId,
		CurrentPrice: startPrice,
		UpdateChannel: utils.MakeDuplicator(),
	}

	stock.PriceChanger = &RandomPrice{
		RunPercent:      runPercent,
		TargetPrice:     100.0,
		PercentToChange: 100,
		Volatility:      5,
	}

	manager.stocks[tickerId] = &stock

	manager.StockUpdateChannel.RegisterInput(stock.UpdateChannel.GetOutput())
	return &stock
}

func buildStockManager() *StockManager{
	return &StockManager{
		stocks: make(map[string]*Stock),
	}
}


//Stock type for storing the stock information
type Stock struct {
	Name         string  `json:"name"`
	TickerId     string  `json:"ticker_id"`
	CurrentPrice float64 `json:"current_price"`
	PriceChanger PriceChange `json:"price_changer"`
	UpdateChannel *utils.ChannelDuplicator
	lock utils.Lock
}


// Some thing that can take in a stock and change the current price
type PriceChange interface {
	change(stock *Stock)
}

// Random Price implements priceChange
type RandomPrice struct {
	RunPercent      float64 `json:"run_percent"`
	TargetPrice     float64 `json:"target_price"`
	PercentToChange float64 `json:"change_percent"`
	Volatility      float64 `json:"volatility"`
}

//change the stock using the changer
func (randPrice *RandomPrice) change(stock *Stock){
	stock.lock.Acquire()
	defer stock.lock.Release()

	if rand.Float64() <= randPrice.RunPercent {
		return
	}
	if rand.Float64() <= randPrice.PercentToChange {
		randPrice.changeValues()
	}

	//can make this a lot more interesting, like adding in the ability for it to drop
	change :=  (randPrice.TargetPrice - stock.CurrentPrice) /
		MapNum(randPrice.Volatility, volatilityMin, volatilityMax, volatilityMinTurns, volatilityMaxTurns)
	fmt.Println("change: ", change)

	stock.CurrentPrice = stock.CurrentPrice + change

	stock.UpdateChannel.Offer(stock)


}

// change the price of the changer
func (randPrice *RandomPrice)changeValues(){

	// get what the upper and lower bounds in % of the current price
	window := MapNum(randPrice.Volatility, volatilityMin, volatilityMax, 0, 0.3)
	fmt.Println("window", window)
	// select a random number on +- that %
	newTarget := MapNum(rand.Float64(), 0, 1, randPrice.TargetPrice * (1 - window/2), randPrice.TargetPrice * (1 + window))

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

