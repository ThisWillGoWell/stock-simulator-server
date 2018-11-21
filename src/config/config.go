package config

import (
	"io/ioutil"
	"os"

	"github.com/stock-simulator-server/src/level"

	"github.com/stock-simulator-server/src/app"
	"github.com/stock-simulator-server/src/items"
	"github.com/stock-simulator-server/src/log"
)

const (
	CONFIG_FOLDER = "CONFIG_FOLDER"
	SEED_JSON     = "SEED_JSON"
	LEVELS_JSON   = "LEVELS_JSON"
	ITEMS_JSON    = "ITEMS_JSON"
)

func Seed() {
	configFolder := os.Getenv(CONFIG_FOLDER)
	if configFolder == "" {
		log.Log.Error("CONFIG_FOLDER not defined")
		return
	}
	seedJson := os.Getenv(SEED_JSON)
	if configFolder == "" {
		log.Log.Error("SEED_JSON is not defined")
		return
	}
	data, err := ioutil.ReadFile(configFolder + seedJson)
	if err != nil {
		log.Log.Error("error reading file", err)
		return
	}
	app.LoadConfig(data)

}

func LoadConfigs() {
	configFolder := os.Getenv(CONFIG_FOLDER)
	if configFolder == "" {
		log.Log.Error("CONFIG_FOLDER not defined")
		return
	}
	levelJson := os.Getenv(LEVELS_JSON)
	if configFolder == "" {
		log.Log.Error("LEVELS_JSON is not defined")
		return
	}
	data, err := ioutil.ReadFile(configFolder + levelJson)
	if err != nil {
		log.Log.Error("error reading file", err)
		return
	}
	level.LoadLevels(data)

	itemsJson := os.Getenv(ITEMS_JSON)
	if configFolder == "" {
		log.Log.Error("ITEMS_JSON is not defined")
		return
	}
	data, err = ioutil.ReadFile(configFolder + itemsJson)
	if err != nil {
		log.Log.Error("error reading file", err)
		return
	}
	items.LoadItemConfig(data)

}
