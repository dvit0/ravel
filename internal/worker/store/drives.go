package store

import (
	"encoding/json"

	"github.com/valyentdev/ravel/pkg/types"
	bolt "go.etcd.io/bbolt"
)

func (store *Store) StoreRavelDrive(id string, driveSpec types.RavelDriveSpec, internal bool) (*types.RavelDrive, error) {
	drive := types.RavelDrive{RavelDriveSpec: &types.RavelDriveSpec{Size: driveSpec.Size}, Id: id, Internal: internal}

	err := store.db.Update(func(tx *bolt.Tx) error {
		json, err := json.Marshal(drive)
		if err != nil {
			return err
		}

		bucket := tx.Bucket([]byte("drives"))

		err = bucket.Put([]byte(drive.Id), json)

		return err
	})
	if err != nil {
		return nil, err
	}

	return &drive, nil
}

func (store *Store) GetRavelDrive(id string) (types.RavelDrive, error) {
	var drive types.RavelDrive

	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("drives"))

		driveBytes := bucket.Get([]byte(id))

		json.Unmarshal(driveBytes, &drive)

		return nil
	})
	if err != nil {
		return types.RavelDrive{}, err
	}

	return drive, nil
}

func (store *Store) ListRavelDrives() ([]types.RavelDrive, error) {
	drivesList := []types.RavelDrive{}

	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("drives"))

		err := bucket.ForEach(func(k, v []byte) error {
			var drive types.RavelDrive

			err := json.Unmarshal(v, &drive)
			if err != nil {
				return err
			}

			drivesList = append(drivesList, drive)

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

	return drivesList, nil
}

func (store *Store) UpdateRavelDrive(id string, driveSpec types.RavelDriveSpec) error {
	drive, err := store.GetRavelDrive(id)
	if err != nil {
		return err
	}

	if driveSpec.Size != 0 {
		drive.Size = driveSpec.Size
	}

	if driveSpec.Name != "" {
		drive.Name = driveSpec.Name
	}

	err = store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("drives"))

		json, err := json.Marshal(drive)
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

func (store *Store) DeleteRavelDrive(id string) error {
	err := store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("drives"))

		err := bucket.Delete([]byte(id))

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
