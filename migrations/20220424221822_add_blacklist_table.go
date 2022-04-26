package migrations

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddBlacklistTable, downAddBlacklistTable)
}

func upAddBlacklistTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS blacklist (
			subnet CIDR PRIMARY KEY
		);
	`)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func downAddBlacklistTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE blacklist;
	`)
	if err != nil {
		return err
	}
	return nil
}
