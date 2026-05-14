package file

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	interval  = 1 * time.Hour
	batchSize = 100
)

func (s *storage) RunCleaner() {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		clean(s.db)
	}()
}

func clean(db *sqlx.DB) {
	for {
		res, err := db.Exec(deleteExpiredQuery(), batchSize)
		if err != nil {
			log.Print(err)
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			log.Print(err)
			return
		}

		if rows == 0 {
			return
		}
	}
}
