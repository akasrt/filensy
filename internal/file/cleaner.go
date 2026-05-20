package file

import (
	"log"
	"time"

	"github.com/akasrt/filensy/internal/database"
	"github.com/akasrt/filensy/internal/filestore"
	"github.com/jmoiron/sqlx"
)

const (
	interval  = 1 * time.Hour
	batchSize = 100
)

func RunCleaner() {
	db := database.GetDB()
	fileStore := filestore.Get()
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		clean(db, fileStore)
	}()
}

func clean(db *sqlx.DB, fileStore filestore.FileStore) {
	for {
		now := time.Now()
		var storageKeys []string
		err := db.Select(&storageKeys, getExpiredStorageKeys(), now, batchSize)
		if err != nil {
			log.Print(err)
		}

		res, err := db.Exec(deleteExpiredQuery(), now, batchSize)
		if err != nil {
			log.Print(err)
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			log.Print(err)
			return
		}

		for _, key := range storageKeys {
			err = fileStore.Delete(key)
			if err != nil {
				log.Print(err)
			}
		}

		if rows == 0 {
			return
		}

	}
}
