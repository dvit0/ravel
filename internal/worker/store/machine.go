package store

import (
	"encoding/json"
	"errors"

	"github.com/valyentdev/ravel/pkg/types"
	bolt "go.etcd.io/bbolt"
)

func (store *Store) StoreRavelMachine(ravelMachine *types.RavelMachine) error {

	err := store.db.Update(func(tx *bolt.Tx) error {
		json, err := json.Marshal(ravelMachine)
		if err != nil {
			return err
		}

		bucket := tx.Bucket([]byte("machines"))

		err = bucket.Put([]byte(ravelMachine.Id), json)

		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (store *Store) GetRavelMachine(id string) (*types.RavelMachine, bool, error) {
	var ravelMachine types.RavelMachine
	found := false
	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("machines"))

		ravelMachineBytes := bucket.Get([]byte(id))
		if ravelMachineBytes == nil {
			return nil
		}
		found = true

		json.Unmarshal(ravelMachineBytes, &ravelMachine)

		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return &ravelMachine, found, nil
}

func (store *Store) ListRavelMachines() ([]types.RavelMachine, error) {
	ravelMachinesList := []types.RavelMachine{}

	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("machines"))
		if bucket == nil {
			return errors.New("ravelMachines bucket not found")
		}

		err := bucket.ForEach(func(k, v []byte) error {
			var ravelMachine types.RavelMachine

			err := json.Unmarshal(v, &ravelMachine)
			if err != nil {
				return err
			}

			ravelMachinesList = append(ravelMachinesList, ravelMachine)

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ravelMachinesList, nil
}

func (store *Store) UpdateRavelMachine(id string, updateRavelMachine func(*types.RavelMachine)) error {
	ravelMachine, found, err := store.GetRavelMachine(id)
	if err != nil {
		return err
	}
	if !found {
		return errors.New("machine not found")
	}

	updateRavelMachine(ravelMachine)

	err = store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("machines"))

		json, err := json.Marshal(ravelMachine)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(id), json)

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (store *Store) DeleteRavelMachine(id string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("machines"))

		err := bucket.Delete([]byte(id))

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
