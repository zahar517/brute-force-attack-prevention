package migrations

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddWhitelistTable, downAddWhitelistTable)
}

func upAddWhitelistTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS whitelist (
			subnet CIDR PRIMARY KEY
		);
	`)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func downAddWhitelistTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE whitelist;
	`)
	if err != nil {
		return err
	}
	return nil
}
