package main
import (
	"stock-server/wallstreet"
	"time"
)

func main(){
	stockManager := wallstreet.NewStockManager()
	stockManager.AddStock("CHUNT", "Conner Hunt's Hats", 420.69)
	stockManager.StartSimulateStocks(3 * time.Second)

}


