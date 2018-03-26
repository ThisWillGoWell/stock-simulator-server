package database

import (
	"github.com/stock-simulator-server/src/valuable"
)

var(
	stocksTableName = `stocks`
	stocksTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + stocksTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`ticker_id text NOT NULL,` +
		`name text NOT NULL,` +
		`current_price numeric(16, 4),` +
		`PRIMARY KEY(uuid)` +
	`);`

	stocksTableUpdateInsert = `INSERT into ` + stocksTableName + `(uuid, ticker_id, name, current_price) values($1, $2, $3, $4) `+
		`ON CONFLICT (uuid) DO UPDATE SET current_price=EXCLUDED.current_price`

	stocksTableQueryStatement = ""
	//getCurrentPrice()
)

func initStocks(){
	tx, err := db.Begin()
	if err != nil{
		db.Close()
		panic("could not begin stocks init: " + err.Error())
	}
	_,err = tx.Exec(stocksTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating metrics table " + err.Error())
	}
	tx.Commit()
}

func runStockUpdate(){
	stockUpdateChannel := valuable.ValuableUpdateChannel.GetBufferedOutput(100)
	go func(){
		for stockUpdated := range stockUpdateChannel{
			stock := stockUpdated.(*valuable.Stock)
			updateStock(stock)
		}
	}();
}

func updateStock(stock *valuable.Stock) {
	dbLock.Acquire("update-stock")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin stocks init")
	}
	_, err = tx.Exec(stocksTableUpdateInsert, stock.Uuid, stock.TickerId, stock.Name, stock.CurrentPrice)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}

func populateStocks(){

}
