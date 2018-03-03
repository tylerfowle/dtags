package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/boltdb/bolt"
	cmdline "github.com/galdor/go-cmdline"
	"github.com/ryanuber/columnize"
)

type database struct {
	bucket     []byte
	dbDir      string
	dbFile     string
	db         string
	currentDir string
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
	cmdline := cmdline.New()

	cmdline.AddCommand("add", "add tag")
	cmdline.AddCommand("del", "delete tag")
	cmdline.AddCommand("get", "get tag")
	cmdline.AddCommand("list", "list all tags")
	cmdline.AddCommand("shell", "list all tags")

	cmdline.Parse(os.Args)

	info := database{}
	info.bucket = []byte("dtags")
	info.dbDir = usr.HomeDir + "/.dtags/go/"
	info.dbFile = "dt.db"
	info.db = info.dbDir + info.dbFile

	info.currentDir = currentDir

	info.args = []string(cmdline.CommandArgumentsValues())
	info.value = []byte(info.currentDir)
	if len(info.args) > 0 {
		info.key = []byte(info.args[0])
	}

	// make the path to the db file if it doesnt exist
	os.MkdirAll(info.dbDir, os.ModePerm)

	switch cmdline.CommandName() {
	case "add":
		addKeyToDatabase(info)
	case "del":
		deleteKeyFromDatabase(info)
	case "get":
		getPathFromTag(info)
	case "list":
		listAllKeysInDatabase(info)
	case "shell":
		shell(info)
	default:
		//  TODO:  <03-03-18, yourname> // default command should search tags and open shell
		shell(info)
	}

}

func addKeyToDatabase(info database) {
	fmt.Printf("adding tag %v to database\n", info.args)

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
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func deleteKeyFromDatabase(info database) {
	fmt.Printf("running command \"del\" with arguments %v\n", info.args)

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
	}

}

func getPathFromTag(info database) string {
	fmt.Printf("getting %v from database\n", info.args)

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

		val = bucket.Get(info.key)
		fmt.Println(string(val))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return string(val)
}

func listAllKeysInDatabase(info database) {
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

		var unformattedList []string
		for k, v := c.First(); k != nil; k, v = c.Next() {
			unformattedList = append(unformattedList, fmt.Sprintf("%s|%s\n", k, v))
		}

		formattedList := columnize.SimpleFormat(unformattedList)
		fmt.Println(formattedList)

		return nil

	})

}

func shell(info database) {

	// Get the current user.
	me, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// setup the path to launch the shell at
	cwd := getPathFromTag(info)

	// Set an environment variable.
	os.Setenv("SOME_VAR", "1")

	// Transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	// Start up a new shell.
	// Note that we supply "login" twice.
	// -fpl means "don't prompt for PW and pass through environment."
	fmt.Print(">> Starting a new interactive shell")
	proc, err := os.StartProcess("/usr/bin/login", []string{"login", "-fpl", me.Username}, &pa)
	if err != nil {
		log.Fatal(err)
	}

	// Wait until user exits the shell
	state, err := proc.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// Keep on keepin' on.
	fmt.Printf("<< Exited shell: %s\n", state.String())
}
