package wallstreet

import (
	"math"
	"math/rand"
	"time"
)

const (
	volatilityMin      = 1
	volatilityMax      = 10
	volatilityMinTurns = 1
	volatilityMaxTurns = 100

)

type StockManager struct {
	stocks map[string]*Stock
	stockUpdateChannel chan *Stock
}


func (manager *StockManager) StartSimulateStocks(intervalLength time.Duration){
	go func() {
		for _ , stock := range manager.stocks{
			stock.priceChanger.change(stock)
			manager.stockUpdateChannel <- stock
		}
	}()


}

func (manager *StockManager) AddStock(tickerId, name string, startPrice float64){
	stock := Stock{
		Name: name,
		TickerId: tickerId,
		CurrentPrice: startPrice,
	}

	stock.priceChanger = &RandomPrice{
		targetPrice: 100.0,
		percentToChange: .1,
		volatility: 5,
	}

	manager.stocks[tickerId] = &stock

	manager.stockUpdateChannel <- &stock
}

func (manager *StockManager) ModifyOpenShares(tickerID string, amount float64){
	 manager.stocks[tickerID].openShares = manager.stocks[tickerID].openShares + amount
	 manager.stockUpdateChannel <- manager.stocks[tickerID]
}

func NewStockManager() *StockManager{
	return &StockManager{
		stocks: make(map[string]*Stock),
		stockUpdateChannel: make(chan *Stock),
	}
}


//Stock type for storing the stock information
type Stock struct {
	Name         string  `json:"name"`
	TickerId     string  `json:"ticker_id"`
	CurrentPrice float64 `json:"current_price"`
	priceChanger  PriceChange
	TotalShares  float64 `json:"total_shares"`
	openShares   float64 `json:"open_shares"`
}



// Some thing that can take in a stock and change the current price
type PriceChange interface {
	change(stock *Stock)
}

// Random Price implements priceChange
type RandomPrice struct {
	targetPrice     float64
	percentToChange float64
	volatility      float64
}

//change the stock using the changer
func (randPrice *RandomPrice) change(stock *Stock){
	if rand.Float64() <= randPrice.percentToChange {
		randPrice.changeValues()
	}
	stock.CurrentPrice = stock.CurrentPrice + ((randPrice.targetPrice - stock.CurrentPrice) /
		MapNum(randPrice.volatility, volatilityMin, volatilityMax, volatilityMinTurns, volatilityMaxTurns))
}

// change the price of the changer
func (randPrice *RandomPrice)changeValues(){
	lowerTarget := RoundPlus(randPrice.targetPrice * (1 - randPrice.volatility)/ 11 * 100, 2)
	upperTarget := RoundPlus(randPrice.targetPrice * (1 + randPrice.volatility)/ 11 * 100, 2)
	randPrice.targetPrice = RoundPlus(RandRange(lowerTarget, upperTarget) / 100.0, 2)
	randPrice.volatility = RandRange(volatilityMin, volatilityMax)
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
	return (value - inMin) * (outMax - outMin) / (inMax - inMin) + outMin
}

