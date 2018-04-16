package database

import ()

var (
	chatHistoryTableName            = `chat_history`
	chatHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + chatHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`message text NULL,` +
		`);`
	chatHistoryTSInit = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE; SELECT create_hypertable('` + chatHistoryTableName + `', 'time');`

	chatHistoryTableUpdateInsert = `INSERT INTO ` + chatHistoryTableName + `(time, uuid, message) values (NOW(), $1, $2)`

	chatHistroyTableQueryStatement = "SELECT * FROM " + chatHistoryTableName + " WHERE uuid"
	//getCurrentPrice()
)

func initChatHistory() {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin chat history init: " + err.Error())
	}
	_, err = tx.Exec(chatHistoryTableCreateStatement)
	if err != nil {

	}
	tx.Commit()
	tx, err = ts.Begin()
	_, err = tx.Exec(chatHistoryTSInit)
	if err != nil {
		// pass on error
	}
	tx.Commit()
}

func SaveChatMessage(uuid, message string) {

	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin portfolio init")
	}
	_, err = tx.Exec(chatHistoryTableUpdateInsert, uuid, message)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert portfolio in table " + err.Error())
	}
	tx.Commit()
}
