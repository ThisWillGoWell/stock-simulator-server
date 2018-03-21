package database

import (
	"github.com/stock-simulator-server/src/valuable"
	"database/sql"
)

var(
	tableName = "stock"
	tableCreateStmt = `CREATE TABLE IF NOT EXISTS ` + tableName +
		`( ` +
		`id serial,` +
		`Name NOT NULL,` +
		`TickerId text NOT NULL,` +
		`CurrentPrice numeric NOT NULL`+
		`PRIMARY KEY (Name)` +
		`);`
	tableUpdatePriceStmt = `INSERT into ` + tableName + `(zone_id, type, TickerId, CurrentPrice) values($1, $2, $3) `+
		`ON CONFLICT (TickerId) DO UPDATE SET CurrentPrice=EXCLUDED.CurrentPrice`
	getCurrentPrice()
)

func startUpdateDatabse(){
	connStr := "postgresql://postgres:5433?user=%sdbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)

}

func updateValueable(){
	update := valuable.ValuableUpdateChannel.GetBufferedOutput(10)


}
