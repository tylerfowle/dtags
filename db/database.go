package db

import (
	"log"
	"os"
	"os/user"

	"github.com/boltdb/bolt"
)

type Database struct {
	Instance         *bolt.DB
	Bucket           []byte
	Datafile         string
	CurrentDirectory string
}

func Init() (error, Database) {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	d := Database{}
	d.Datafile = u.HomeDir + "/.dtags/go/dt.db"
	d.Bucket = []byte("dtags")
	d.CurrentDirectory, _ = os.Getwd()

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	d.Instance, err = bolt.Open(d.Datafile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = d.Instance.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(d.Bucket)
		return err
	})

	return err, d
}

func (d *Database) AddKey(k string, v string) error {
	return d.Instance.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(d.Bucket).Put([]byte(k), []byte(v))
	})
}

func (d *Database) DeleteKey(k string) error {
	return d.Instance.Update(func (tx *bolt.Tx) error {
		return tx.Bucket(d.Bucket).Delete([]byte(k))
	})
}

func (d *Database) GetTags() []string {
	var tags []string
	for key := range d.All() {
		tags = append(tags, key)
	}

	return tags
}

func (d *Database) GetValue(k string) string {
	var val []byte
	_ = d.Instance.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(d.Bucket).Get([]byte(k))

		return nil
	})

	return string(val)
}

func (d *Database) All() map[string]string {
	var m map[string]string
	_ = d.Instance.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(d.Bucket).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			m[string(k)] = string(v)
		}

		return nil
	})

	return m
}

