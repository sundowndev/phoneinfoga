package cache

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/sundowndev/phoneinfoga.v2/scanners"
	"time"
)

const (
	defaultFilename = "cache.db"
	scanTasksBucket = "scans"
	resultsBucket   = "numbers"
)

// DefaultKeyExpiration is the default value for temporary data (7 days)
const DefaultKeyExpiration = 7 * (24 * time.Hour)

type Cache struct {
	Directory string
	Filename  string
	db        *bolt.DB
}

func (c *Cache) Init() error {
	var filename string

	if c.db != nil {
		return fmt.Errorf("database was already initialized")
	}

	if c.Filename == "" {
		filename = defaultFilename
	}

	db, err := bolt.Open(fmt.Sprintf("%s/%s", c.Directory, filename), 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	c.db = db

	err = c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(scanTasksBucket))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(resultsBucket))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) CacheResults(n *scanners.ScanResult) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(resultsBucket))

		data, err := json.Marshal(n)
		if err != nil {
			return err
		}

		err = b.Put([]byte(n.Local.International), data)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetResults(n *scanners.Number) (*scanners.ScanResult, error) {
	results := &scanners.ScanResult{}

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(resultsBucket))

		data := b.Get([]byte(n.International))

		err := json.Unmarshal(data, results)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return results, err
	}

	return results, nil
}

func (c *Cache) Close() error {
	return c.db.Close()
}
