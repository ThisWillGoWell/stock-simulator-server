package app

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/items"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/ledger"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/record"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/user"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"
)

func LoadFromDb() error {
	if err := loadStocks(); err != nil {
		return err
	}
	if err := loadUsers(); err != nil {
		return err
	}
	if err := loadPortfolio(); err != nil {
		return err
	}

	if err := loadEffects(); err != nil {
		return err
	}
	if err := loadItems(); err != nil {
		return err
	}
	if err := loadLedgers(); err != nil {
		return err
	}

	if err := loadRecords(); err != nil {
		return err
	}

	if err := loadNotificiaons(); err != nil {
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

func loadRecords() error {
	models, err := database.Db.GetRecords()
	if err != nil {
		return err
	}
	for _, m := range models {
		record.MakeRecord(m, false)
	}
	return nil
}

func loadUsers() error {
	models, err := database.Db.GetUsers()
	if err != nil {
		return err
	}
	for _, m := range models {
		if _, err := user.MakeUser(m); err != nil {
			return err
		}
	}
	return nil
}

func loadEffects() error {
	effects, err := database.Db.GetEffects()
	if err != nil {
		return err
	}
	for _, m := range effects {
		if _, err := effect.MakeEffect(m, false); err != nil {
			return err
		}
	}
	return nil
}

func loadItems() error {
	models, err := database.Db.GetItems()
	if err != nil {
		return err
	}
	for _, m := range models {
		if _, err := items.MakeItem(m); err != nil {
			return err
		}
	}
	return nil
}

func loadLedgers() error {
	models, err := database.Db.GetLedgers()
	if err != nil {
		return err
	}
	for _, m := range models {
		if _, err := ledger.MakeLedgerEntry(m, false); err != nil {
			return err
		}
	}
	return nil
}

func loadNotificiaons() error {
	models, err := database.Db.GetNotification()
	if err != nil {
		return err
	}
	for _, m := range models {
		notification.MakeNotification(m)
	}
	return nil
}

func loadPortfolio() error {
	models, err := database.Db.GetPortfolios()
	if err != nil {
		return err
	}
	for _, m := range models {
		if _, err := portfolio.MakePortfolio(m, false); err != nil {
			return err
		}
	}
	return nil
}

func loadStocks() error {
	models, err := database.Db.GetStocks()
	if err != nil {
		return err
	}

	for _, m := range models {
		if _, err := valuable.MakeStock(m); err != nil {
			return err
		}
	}
	return nil
}
