package objects

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/items"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/ledger"
)

func LoadFromDb() error {
	if err := loadEffects(); err != nil {
		return err
	}
	if err := loadItems(); err != nil {
		return err
	}
	if err := loadLedgers(); err != nil {
		return err
	}
	if err := loadNotificiaons(); err != nil {
		return err
	}
	if err := loadPortfolio(); err != nil {
		return err
	}

	if err := loadStocks(); err != nil {
		return err
	}

	for _, l := range ledger.Entries {
		port := portfolio.Portfolios[l.PortfolioId]
		stock := valuable.Stocks[l.StockId]
		port.UpdateInput.RegisterInput(stock.UpdateChannel.GetBufferedOutput(100))
		port.UpdateInput.RegisterInput(l.UpdateChannel.GetBufferedOutput(100))
	}
	for _, port := range portfolio.Portfolios {
		port.Update()
	}

	return nil

}

func loadEffects() error {
	effects,err  := database.Db.GetEffects()
	if err != nil {
		return err
	}
	for uuid, m := range effects{
		effect.MakeEffect(uuid, m, false )
	}
	return nil
}

func loadItems() error {
	models,err  := database.Db.GetItems()
	if err != nil {
		return err
	}
	for _, m := range models{
		items.MakeItem(m)
	}
	return nil
}

func loadLedgers() error {
	models,err  := database.Db.GetLedgers()
	if err != nil {
		return err
	}
	for _, m := range models{
		ledger.MakeLedgerEntry(m, false)
	}
	return nil
}

func loadNotificiaons() error {
	models,err  := database.Db.GetNotification()
	if err != nil {
		return err
	}
	for _, m := range models{
		notification.MakeNotification(m)
	}
	return nil
}


func loadPortfolio() error {
	models,err  := database.Db.GetPortfolios()
	if err != nil {
		return err
	}
	for _, m := range models{
		portfolio.MakePortfolio(m, false)
	}
	return nil
}

func loadStocks() error {
	models,err  := database.Db.GetStocks()
	if err != nil {
		return err
	}

	for _, m := range models{
		valuable.MakeStock(m)
	}
	return nil
}

