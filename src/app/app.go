package app

import (
	"fmt"
	"github.com/ThisWillGoWell/stock-simulator-server/src/app/aws"
	"github.com/ThisWillGoWell/stock-simulator-server/src/app/config"
	"github.com/ThisWillGoWell/stock-simulator-server/src/app/seed"
	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/ThisWillGoWell/stock-simulator-server/src/game/level"
	"github.com/ThisWillGoWell/stock-simulator-server/src/game/order"
	"github.com/ThisWillGoWell/stock-simulator-server/src/id/change"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/items"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/web/http"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/sender"
	"io/ioutil"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/web/session"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"
)

func LoadConfigs() {

}
func App() {

	c := config.FromEnv()
	secret, err := aws.GetDatabaseSecret(c.Environment)
	if err != nil {
		panic(err)
	}

	if err := database.InitDatabase(c.EnableDb, c.EnableDbWrites, secret.Host, fmt.Sprintf("%d", secret.Port), secret.Username, secret.Password, "postgres"); err != nil {
		panic(err)
	}
	log.Log.Info("Starting App")
	if c.EnableDb {
		log.Log.Info("loading from database")
		if err := LoadFromDb(); err != nil {
			panic(err)
		}
		log.Log.Info("done")
	}

	portfolio.UpdateAll()
	//valuable.ValuablesLock.EnableDebug()
	//ledger.EntriesLock.EnableDebug()

	//discordAlertToken := secrest.DiscordToken
	//
	//var alertWriter io.Writer
	//if discordAlertToken != "" {
	//	alertWriter = alert.Init(discordAlertToken, "504397270075179029")
	//} else {
	//	// if there is discord token, discard all alerts
	//	alertWriter = ioutil.Discard
	//}
	//log.Init(alertWriter)
	//log.Alerts.Info("Starting App")
	log.Log.Info("Connecting wires")
	//Wiring of system
	wires.ConnectWires()
	log.Log.Info("done")

	//this takes the subscribe output and converts it to a message
	change.StartDetectChanges()
	session.StartSessionCleaner()
	sender.RunGlobalSender()
	effect.RunEffectCleaner()

	order.Run()
	valuable.StartStockStimulation()

	// read in the levels and items config
	data, err := ioutil.ReadFile(c.LevelsJson)
	if err != nil {
		log.Log.Error("error reading file", err)
		return
	}
	level.LoadLevels(data)

	data, err = ioutil.ReadFile(c.ItemsJson)
	if err != nil {
		log.Log.Error("error reading file", err)
		return
	}
	items.LoadItemConfig(data)

	if c.SeedObjects {
		log.Log.Info("Seeding From Config")
		seed.SeedObjects(c.ObjectsJson)
		log.Log.Info("done")
	}
	log.Log.Info("Starting Handlers")
	http.StartHandlers()
}
