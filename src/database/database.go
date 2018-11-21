package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/stock-simulator-server/src/effect"

	"github.com/stock-simulator-server/src/log"

	"github.com/stock-simulator-server/src/record"

	"github.com/stock-simulator-server/src/items"

	"github.com/stock-simulator-server/src/notification"

	_ "github.com/lib/pq"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
	"github.com/stock-simulator-server/src/wires"
)

var db *sql.DB
var ts *sql.DB

var dbLock = lock.NewLock("db lock")

func InitDatabase(disableDbWrite bool) {
	dbConStr := os.Getenv("DB_URI")
	// if the env is not set, default to use the local host default port
	database, err := sql.Open("postgres", dbConStr)
	fmt.Println(dbConStr)
	if err != nil {
		log.Alerts.Fatal("could not connect to database: " + err.Error())
		log.Log.Fatal("could not connect to database: " + err.Error())
		panic("could not connect to database: " + err.Error())
	}
	db = database

	for i := 0; i < 10; i++ {
		err := db.Ping()

		if err == nil {
			break
		}
		fmt.Println("	waitng for connection to db")
		<-time.After(time.Second)
	}

	initLedger()
	initStocks()
	initPortfolio()
	initStocksHistory()
	initPortfolioHistory()
	initNotification()
	initItems()
	initLedgerHistory()
	initAccount()
	initRecordHistory()
	initEffect()

	populateUsers()
	populateStocks()
	populatePortfolios()
	populateLedger()
	populateItems()
	populateNotification()
	populateRecords()
	populateEffects()

	for _, l := range ledger.Entries {
		port := portfolio.Portfolios[l.PortfolioId]
		stock := valuable.Stocks[l.StockId]
		port.UpdateInput.RegisterInput(stock.UpdateChannel.GetBufferedOutput(100))
		port.UpdateInput.RegisterInput(l.UpdateChannel.GetBufferedOutput(100))
	}
	for _, port := range portfolio.Portfolios {
		port.Update()
	}

	runHistoricalQueries()
	if !disableDbWrite {
		fmt.Println("starting db writer")
		go databaseWriter()
	}

}

func databaseWriter() {
	go func() {
		portfolioDBWrite := duplicator.MakeDuplicator("portfolio-db-write")
		portfolioDBWrite.RegisterInput(wires.PortfolioNewObject.GetBufferedOutput(100))
		portfolioDBWrite.RegisterInput(wires.PortfolioUpdate.GetBufferedOutput(100))
		write := portfolioDBWrite.GetBufferedOutput(1000)
		for val := range write {
			writePortfolio(val.(*portfolio.Portfolio))
			writePortfolioHistory(val.(*portfolio.Portfolio))
		}
	}()

	go func() {
		userDBWrite := duplicator.MakeDuplicator("user-db-write")
		userDBWrite.RegisterInput(wires.UsersNewObject.GetBufferedOutput(100))
		userDBWrite.RegisterInput(wires.UsersUpdate.GetBufferedOutput(100))
		write := userDBWrite.GetBufferedOutput(1000)
		for val := range write {
			writeUser(val.(*account.User))
		}
	}()

	go func() {
		ledgerDBWrite := duplicator.MakeDuplicator("ledger-db-write")
		ledgerDBWrite.RegisterInput(wires.LedgerNewObject.GetBufferedOutput(100))
		ledgerDBWrite.RegisterInput(wires.LedgerUpdate.GetBufferedOutput(100))
		write := ledgerDBWrite.GetBufferedOutput(1000)
		for val := range write {
			writeLedger(val.(*ledger.Entry))
			writeLedgerHistory(val.(*ledger.Entry))
		}
	}()

	go func() {
		itemsDBWrite := duplicator.MakeDuplicator("items-db-write")
		itemsDBWrite.RegisterInput(wires.ItemsNewObjects.GetBufferedOutput(100))
		itemsDBWrite.RegisterInput(wires.ItemsUpdate.GetBufferedOutput(100))
		write := itemsDBWrite.GetBufferedOutput(1000)
		for val := range write {
			writeItem(val.(*items.Item))
		}
	}()
	go func() {
		itemDelete := wires.ItemsDelete.GetBufferedOutput(100)
		for item := range itemDelete {
			deleteItem(item.(*items.Item))
		}
	}()

	go func() {
		notificationDBWrite := duplicator.MakeDuplicator("notification-db-write")
		notificationDBWrite.RegisterInput(wires.NotificationNewObject.GetBufferedOutput(100))
		notificationDBWrite.RegisterInput(wires.NotificationUpdate.GetBufferedOutput(100))
		write := notificationDBWrite.GetBufferedOutput(1000)
		for val := range write {
			writeNotification(val.(*notification.Notification))
		}
	}()
	go func() {
		notificationsDelete := wires.NotificationsDelete.GetBufferedOutput(100)
		for note := range notificationsDelete {
			deleteNotification(note.(*notification.Notification))
		}
	}()

	go func() {
		stockDBWrite := duplicator.MakeDuplicator("stock-9-write")
		stockDBWrite.RegisterInput(wires.StocksNewObject.GetBufferedOutput(100))
		stockDBWrite.RegisterInput(wires.StocksUpdate.GetBufferedOutput(100))
		write := stockDBWrite.GetBufferedOutput(1000)
		for val := range write {
			writeStock(val.(*valuable.Stock))
			writeStockHistory(val.(*valuable.Stock))
		}
	}()

	go func() {
		recordsWrite := wires.RecordsNewObject.GetBufferedOutput(1000)
		for val := range recordsWrite {
			writeRecordHistory(val.(*record.Record))
		}
	}()

	go func() {
		effectDBWrite := duplicator.MakeDuplicator("-db-write")
		effectDBWrite.RegisterInput(wires.EffectsNewObject.GetBufferedOutput(100))
		effectDBWrite.RegisterInput(wires.EffectsUpdate.GetBufferedOutput(100))
		write := effectDBWrite.GetBufferedOutput(1000)
		for val := range write {
			writeEffect(val.(*effect.Effect))
		}
	}()
	go func() {
		effectDelete := wires.EffectsDelete.GetBufferedOutput(100)
		for val := range effectDelete {
			deleteEffect(val.(*effect.Effect))
		}
	}()

}
