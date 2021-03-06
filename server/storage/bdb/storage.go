package bdb

import (
	"github.com/asdine/storm"
	"github.com/fireeye/gocrack/server/storage"
)

// CurrentStorageVersion describes what storage version we're on and can be used to determine
// if the underlying database structure needs to change
const CurrentStorageVersion = 1.0

var bucketInternalConfig = []byte("config")

func init() {
	storage.Register("bdb", &Driver{})
}

type Driver struct{}

func (s *Driver) Open(cfg storage.Config) (storage.Backend, error) {
	return Init(cfg)
}

// BoltBackend is a storage backend for GoCrack built ontop of Bolt/LMDB
type BoltBackend struct {
	// contains filtered or unexported fields
	db *storm.DB
}

// Init creates a database instance backed by boltdb
func Init(cfg storage.Config) (storage.Backend, error) {
	db, err := storm.Open(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	bb := &BoltBackend{db: db}
	if err = bb.checkSchema(); err != nil {
		return nil, err
	}
	return bb, nil
}

// Close the resources used by bolt
func (s *BoltBackend) Close() error {
	return s.db.Close()
}

func (s *BoltBackend) checkSchema() error {
	var opt interface{}
	// check our internal version
	if err := s.db.Get("options", "version", &opt); err != nil {
		if err == storm.ErrNotFound {
			if err := s.db.Set("options", "version", CurrentStorageVersion); err != nil {
				return err
			}
			opt = CurrentStorageVersion
			err = nil
		}
	}
	// XXX(cschmitt): What do if version changes?
	return nil
}
