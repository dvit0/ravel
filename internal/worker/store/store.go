package store

import (
	"github.com/charmbracelet/log"
	bolt "go.etcd.io/bbolt"
)

type Store struct {
	db *bolt.DB
}

func initializeBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("drives"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("machines"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("drivers"))
		if err != nil {
			return err
		}

		return nil
	},
	)

}

func NewStore() Store {
	db, err := bolt.Open("/var/lib/ravel/worker.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = initializeBuckets(db)

	if err != nil {
		log.Fatal(err)
	}

	return Store{
		db: db,
	}
}

func (store *Store) CloseStore() error {
	return store.db.Close()
}
