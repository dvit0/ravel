package store

import (
	"encoding/json"

	"github.com/hashicorp/go-plugin"
	bolt "go.etcd.io/bbolt"
)

type RavelDriverReattachConfig struct {
	Protocol        plugin.Protocol
	ProtocolVersion int
	Addr            DriverAddr
	Pid             int
}

type DriverAddr struct {
	Net     string
	Address string
}

func (addr *DriverAddr) Network() string {
	return addr.Net
}

func (addr *DriverAddr) String() string {
	return addr.Address
}

func (store *Store) GetDriverReattachConfig(name string) (*plugin.ReattachConfig, error) {
	var drive RavelDriverReattachConfig
	found := true

	err := store.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("drivers"))

		driveBytes := bucket.Get([]byte(name))
		if driveBytes == nil {
			found = false
			return nil
		}

		err := json.Unmarshal(driveBytes, &drive)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &plugin.ReattachConfig{
		Protocol:        drive.Protocol,
		ProtocolVersion: drive.ProtocolVersion,
		Addr:            &drive.Addr,
		Pid:             drive.Pid,
	}, nil
}

func (store *Store) SetDriverReattachConfig(name string, config *plugin.ReattachConfig) error {
	ravelDriverReattachConfig := RavelDriverReattachConfig{
		Protocol:        config.Protocol,
		ProtocolVersion: config.ProtocolVersion,
		Addr:            DriverAddr{Net: config.Addr.Network(), Address: config.Addr.String()},
		Pid:             config.Pid,
	}

	err := store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("drivers"))

		json, err := json.Marshal(ravelDriverReattachConfig)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(name), json)

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
