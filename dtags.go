package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/ryanuber/columnize"
)

type database struct {
	bucket     []byte
	dbDir      string
	dbFile     string
	db         string
	currentDir string
	subcommand string
	key        []byte
	value      []byte
	args       []string
}

func main() {

	// user info
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// ##################################################
	// Command Line
	// ##################################################

	// setup db vars
	info := database{}
	info.bucket = []byte("dtags")
	info.dbDir = usr.HomeDir + "/.dtags/go/"
	info.dbFile = "dt.db"
	info.db = info.dbDir + info.dbFile

	// current directory
	info.currentDir = currentDir
	info.value = []byte(info.currentDir)

	if len(os.Args[0:]) < 1 {
		fmt.Printf("no options")
		os.Exit(1)
	}

	// arguments
	info.subcommand = string(os.Args[1])
	info.args = []string(os.Args[2:])

	// make the path to the db file if it doesnt exist
	os.MkdirAll(info.dbDir, os.ModePerm)

	switch info.subcommand {
	case "add":
		info.key = []byte(info.args[0])
		if info.key == nil {
			fmt.Printf("not enough arguments")
			os.Exit(1)
		}
		addKeyToDatabase(info)
	case "del":
		info.key = []byte(info.args[0])
		deleteKeyFromDatabase(info)
	case "list":
		listTags(info)
	case "all", "more":
		listAll(info)
	case "shell":
		info.key = []byte(info.args[0])
		shell(info)
	case "completion":
		tagCompletion(info)
	default:
		info.key = []byte(info.subcommand)
		info.args = []string(os.Args[1:])
		shell(info)
	}

}

func addKeyToDatabase(info database) {

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

		err = bucket.Put(info.key, info.value)
		if err == nil {
			fmt.Printf("added tag [%v] with path [%v] to database\n", string(info.key), string(info.value))
		} else {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func deleteKeyFromDatabase(info database) {

	db, err := bolt.Open(info.db, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Delete the key in a different write transaction.
	if err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(info.bucket)).Delete([]byte(info.key))
	}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("tag [%v] deleted\n", string(info.key))
	}

}

func listTags(info database) {

	// return keys on current dir
	db, err := bolt.Open(info.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(info.bucket))
		c := b.Cursor()

		var tags []string
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == info.currentDir {
				tags = append(tags, fmt.Sprintf("%s", k))
			}
		}

		if len(tags) > 0 {
			fmt.Println(tags)
		}

		return nil

	})

}

func getPathFromTag(info database) string {
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

		val = bucket.Get(info.key)
		if val == nil {
			fmt.Printf("no tag %v found\n", info.args)
			os.Exit(1)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return string(val)
}

func listAll(info database) {

	db, err := bolt.Open(info.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(info.bucket))
		c := b.Cursor()

		var unformattedList []string
		for k, v := c.First(); k != nil; k, v = c.Next() {
			unformattedList = append(unformattedList, fmt.Sprintf("%s|%s\n", k, v))
		}

		formattedList := columnize.SimpleFormat(unformattedList)
		// print out the column formatted list
		fmt.Println(formattedList)

		return nil

	})

}

func tagCompletion(info database) {

	db, err := bolt.Open(info.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(info.bucket))
		c := b.Cursor()

		var unformattedList []string
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			unformattedList = append(unformattedList, fmt.Sprintf("%s\n", k))
		}

		formattedList := columnize.SimpleFormat(unformattedList)
		// print out the column formatted list
		fmt.Println(formattedList)

		return nil

	})

}

func shell(info database) {

	// setup the path to launch the shell at
	cwd := getPathFromTag(info)
	if cwd == "" {
		// exit if path is nil
		fmt.Printf("tag not found")
		os.Exit(1)
	}

	// Set an environment variable.
	os.Setenv("DTAGSPID", strconv.Itoa(os.Getpid()))

	fmt.Fprint(os.Stdout, cwd)
	os.Exit(1)
}
