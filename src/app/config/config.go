package config

import (
	"os"
)

type EnvConfig struct {
	Environment string
	EnableDbWrites bool
	EnableDb bool
	SeedObjects bool
	ItemsJson string
	LevelsJson string
	ObjectsJson string
}

func FromEnv() EnvConfig{
	return  EnvConfig{
		Environment:    os.Getenv("ENV"),
		EnableDbWrites: os.Getenv("ENABLE_DB_WRITES") == "true",
		EnableDb:       os.Getenv("ENABLE_DB") == "true",
		ItemsJson: os.Getenv("ITEMS_JSON"),
		LevelsJson: os.Getenv("LEVELS_JSON"),
		SeedObjects: os.Getenv("SEED_OBJECTS") == "true",
		ObjectsJson: os.Getenv("OBJECTS_JSON"),
	}
}




