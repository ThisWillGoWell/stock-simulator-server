package ledger


import "github.com/ThisWillGoWell/stock-simulator-server/src/database"

func LoadEffects() error {
	ledgers,err  := database.Db.GetLedgers()
	if err != nil {
		return err
	}
	for _, m := range ledgers{
		MakeItem(m)
	}
	return nil
}
