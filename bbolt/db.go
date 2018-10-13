package db

import (
	_ "flag"
	"log"

	bolt "github.com/coreos/bbolt"
)

// DB Connection
type BoltDB struct {
	db *bolt.DB
}

// TODO: Fix it
// https://stackoverflow.com/questions/26537806/how-to-access-flags-outside-of-main-package
// dbName variable
var dbName = "gopxe.db"

// DB initilization amd creatin database buckets
func init() {
	var p BoltDB
	err := p.CreateBucket("bootactions")
	if err != nil {
		panic(err)
	}
	err = p.CreateBucket("pxe")
	if err != nil {
		panic(err)
	}
}

// Get DB sessions
func GetSession(dbName string) *bolt.DB {
	// Cretes the DB if it doens't exist
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		panic(err)
	}
	return db
}

// Creating bucket
func (b *BoltDB) CreateBucket(name string) error {
	b.db = GetSession(dbName)
	defer b.db.Close()

	// Start a writable transaction.
	tx, err := b.db.Begin(true)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucket([]byte(name))
	if err != nil {
		log.Printf("Bucket %s already exists\n", name)
	} else {
		log.Printf("Bucket %s created\n", name)
	}

	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return nil
}

// Storing boot action into the database
func (b *BoltDB) PutBootAction(bucket, key, value string) error {
	b.db = GetSession(dbName)
	defer b.db.Close()

	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return nil
}

// Getting bootaction from database
func (b *BoltDB) GetBootAction(bucket, key string) (error, string) {
	b.db = GetSession(dbName)
	defer b.db.Close()
	var vv string
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		//fmt.Printf("Key %s, %s\n", key, v)
		vv = string(v)
		//fmt.Printf("Key1 %s, %s\n", key, vv)
		return nil
	})
	return nil, vv
}

// GetAllBootActions gets all boot action stored in the database
func (b *BoltDB) GetAllBootActions(bucket string) (error, map[string]string) {
	b.db = GetSession(dbName)
	defer b.db.Close()

	items := make(map[string]string)
	b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		a := tx.Bucket([]byte(bucket))

		c := a.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			//fmt.Printf("key=%s, value=%s\n", k, v)
			items[string(k)] = string(v)
		}

		return nil
	})
	return nil, items
}
