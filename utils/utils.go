package utils

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
)

func NewDB(server string, port string, database string, creds ...string) (*sql.DB, error) {
	var (
		credstring string
	)
	if len(creds) > 0 {
		if len(creds) != 2 {
			return nil, fmt.Errorf("not enough credential arguments supplied: %d", len(creds))
		}
		credstring = fmt.Sprintf(";User Id=%s;Password=%s", creds[0], creds[1])
	}
	return sql.Open("sqlserver", fmt.Sprintf("Server=%s;Port=%s;Database=%s%s", server, port, database, credstring))
}

func QueryReturningRowCount(db *sql.Tx, qry string) (int64, error) {
	var i int64
	row := db.QueryRow(qry)
	if err := row.Scan(&i); err != nil {
		return -1, err
	}
	return i, nil
}

type multiwriteCloser struct {
	f io.WriteCloser
	w io.Writer
}

func (mwc *multiwriteCloser) Write(p []byte) (int, error) {
	return mwc.w.Write(p)
}

func (mwc *multiwriteCloser) Close() error {
	return mwc.f.Close()
}

func NewLogWriter(fname string) (io.WriteCloser, error) {
	lf, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	return &multiwriteCloser{
		f: lf,
		w: io.MultiWriter(lf, os.Stdout),
	}, nil
}

func HandleTxFunc(tx *sql.Tx, ok *bool) {
	if !*ok {
		log.Println("rolling back...")
		err := tx.Rollback()
		if err != nil {
			log.Printf("error rolling back tx: %v\n", err)
			return
		}
		log.Println("successfully rolled back tx")
		return
	}
	log.Println("committing...")
	if err := tx.Commit(); err != nil {
		log.Printf("err committing tx: %v\n", err)
	}
}
