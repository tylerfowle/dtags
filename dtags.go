package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/boltdb/bolt"
	cmdline "github.com/galdor/go-cmdline"
)

type database struct {
	bucket []byte
	dbDir  string
	dbFile string
	db     string
}

func main() {

	// user info
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	info := database{}
	info.bucket = []byte("dtags")
	info.dbDir = usr.HomeDir + "/.dtags/go/"
	info.dbFile = "dt.db"
	info.db = info.dbDir + info.dbFile

	// make the path to the db file
	os.MkdirAll(info.db, os.ModePerm)

	// ##################################################
	// Command Line
	// ##################################################
	cmdline := cmdline.New()

	cmdline.AddCommand("add", "add tag")
	cmdline.AddCommand("del", "delete tag")
	cmdline.AddCommand("get", "get tag")
	cmdline.AddCommand("list", "list all tags")

	cmdline.Parse(os.Args)

	switch cmdline.CommandName() {
	case "add":
		addKeyToDatabase(cmdline.CommandArgumentsValues(), info)
	case "del":
		deleteKeyFromDatabase(cmdline.CommandArgumentsValues(), info)
	case "get":
		getKeyFromDatabase(cmdline.CommandArgumentsValues(), info)
	case "list":
		listAllKeysInDatabase(cmdline.CommandArgumentsValues(), info)
	}

}

func addKeyToDatabase(args []string, info database) {
	fmt.Printf("adding tag %v to database\n", args)

	// tag
	key := []byte(args[0])
	// current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	value := []byte(currentDir)

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(info.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// store some data
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(info.bucket)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func deleteKeyFromDatabase(args []string, info database) {
	fmt.Printf("running command \"del\" with arguments %v\n", args)
}

func getKeyFromDatabase(args []string, info database) string {
	fmt.Printf("getting %v from database\n", args)

	// tag
	key := []byte(args[0])

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(info.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var val []byte
	// retrieve the data
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(info.bucket)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found! ", info.bucket)
		}

		val = bucket.Get(key)
		fmt.Println(string(val))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return string(val)
}

func listAllKeysInDatabase(args []string, info database) {
	fmt.Printf("listing all keys in database\n")

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(info.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(info.bucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil

	})

}
