package mssql

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

func Connect(cnstr string) (*sql.Tx, error) {
	var err error
	db, err = sql.Open("sqlserver", cnstr)
	if err != nil {
		return nil, err
	}
	return db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSnapshot})
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func wrapQuery(q string) string {
	qwrap := `BEGIN TRY
%s
END TRY
BEGIN CATCH
THROW;
END CATCH`
	return fmt.Sprintf(qwrap, q)
}
