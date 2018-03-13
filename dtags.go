package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

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

var (
	db   *bolt.DB
	err  error
	info = database{}
)

func init() {
	// setup db vars
	info.bucket = []byte("dtags")
	info.dbDir = getUser().HomeDir + "/.dtags/go/"
	info.dbFile = "dt.db"
	info.db = info.dbDir + info.dbFile

	// setup current directory and value
	info.currentDir = getCurrentDir()
	info.value = []byte(info.currentDir)

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err = bolt.Open(info.db, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func getUser() *user.User {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr
}

func getCurrentDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}

func main() {

	// defer closing the database
	defer db.Close()

	// arguments
	info.subcommand = string(os.Args[1])
	info.args = []string(os.Args[2:])

	// make the path to the db file if it doesnt exist
	os.MkdirAll(info.dbDir, os.ModePerm)

	switch info.subcommand {
	case "add":
		info.key = []byte(strings.ToLower(info.args[0]))

		if len(info.args[0:]) > 1 {
			info.value = []byte(info.args[1])
		}

		addKeyToDatabase(info)
	case "del":
		info.key = []byte(strings.ToLower(info.args[0]))
		deleteKeyFromDatabase(info)
	case "list":
		listTags(info)
	case "ls":
		listAll(info)
	case "completion":
		tagCompletion(info)
	default:
		info.key = []byte(info.subcommand)
		info.args = []string(os.Args[1:])
		shell(info)
	}

}

//func addKeyToDatabase(info database) {
//
//	// if getPathFromTag(info) != "" {
//	// fmt.Printf("Overwrite existing tag? [%s] (y/n)", info.key)
//	// if confirmation() == false {
//	// fmt.Println("action cancelled. no tag added")
//	// os.Exit(0)
//	// }
//	// } else {
//	// fmt.Printf("other")
//	// }
//
//	err = db.Update(func(tx *bolt.Tx) error {
//		bucket, err := tx.CreateBucketIfNotExists(info.bucket)
//		if err != nil {
//			return err
//		}
//
//		err = bucket.Put(info.key, info.value)
//		if err == nil {
//			fmt.Printf("added tag [%v] with path [%v]\n", string(info.key), string(info.value))
//		} else {
//			return err
//		}
//		return nil
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}

//func deleteKeyFromDatabase(info database) {
//
//	if err := db.Update(func(tx *bolt.Tx) error {
//		return tx.Bucket([]byte(info.bucket)).Delete([]byte(info.key))
//	}); err != nil {
//		log.Fatal(err)
//	} else {
//		fmt.Printf("deleted tag [%v]\n", string(info.key))
//	}
//
//}
//
//func listTags(info database) {
//
//	err = db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(info.bucket))
//		c := b.Cursor()
//
//		var tags []string
//		for k, v := c.First(); k != nil; k, v = c.Next() {
//			if string(v) == info.currentDir {
//				tags = append(tags, fmt.Sprintf("%s", k))
//			}
//		}
//
//		if len(tags) > 0 {
//			fmt.Println(tags)
//		}
//
//		return nil
//
//	})
//
//}
//
//func getPathFromTag(info database) string {
//
//	var val []byte
//	// retrieve the data
//	err = db.View(func(tx *bolt.Tx) error {
//		bucket := tx.Bucket(info.bucket)
//		if bucket == nil {
//			return fmt.Errorf("bucket %q not found! ", info.bucket)
//		}
//
//		val = bucket.Get(info.key)
//		if val == nil {
//			fmt.Printf("no tag %v found\n", info.args)
//			os.Exit(1)
//		}
//		return nil
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return string(val)
//}
//
//func listAll(info database) {
//
//	err = db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(info.bucket))
//		c := b.Cursor()
//
//		var unformattedList []string
//		for k, v := c.First(); k != nil; k, v = c.Next() {
//			unformattedList = append(unformattedList, fmt.Sprintf("%s|%s\n", k, v))
//		}
//
//		formattedList := columnize.SimpleFormat(unformattedList)
//		// print out the column formatted list
//		fmt.Println(formattedList)
//
//		return nil
//
//	})
//
//}
//
//func tagCompletion(info database) {
//
//	err = db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(info.bucket))
//		c := b.Cursor()
//
//		var unformattedList []string
//		for k, _ := c.First(); k != nil; k, _ = c.Next() {
//			unformattedList = append(unformattedList, fmt.Sprintf("%s\n", k))
//		}
//
//		formattedList := columnize.SimpleFormat(unformattedList)
//		// print out the column formatted list
//		fmt.Println(formattedList)
//
//		return nil
//
//	})
//
//}
//
//func shell(info database) {
//
//	// setup the path to launch the shell at
//	cwd := getPathFromTag(info)
//	if cwd == "" {
//		// exit if path is nil
//		fmt.Printf("tag not found")
//		os.Exit(1)
//	}
//
//	// Set an environment variable.
//	// os.Setenv("DTAGSPID", strconv.Itoa(os.Getpid()))
//
//	fmt.Fprint(os.Stdout, cwd)
//	os.Exit(1)
//}
//
//func confirmation() bool {
//	var response string
//
//	_, err := fmt.Scanln(&response)
//	if err != nil {
//		log.Fatal(err)
//	}
//	y := []string{"y", "Y", "yes", "Yes", "YES"}
//	n := []string{"n", "N", "no", "No", "NO"}
//
//	response = strings.TrimSpace(response)
//	response = strings.ToLower(response)
//
//	if containsString(y, response) {
//		return true
//	} else if containsString(n, response) {
//		return false
//	} else {
//		fmt.Println("yes or no required:")
//		return confirmation()
//	}
//}
//
//func posString(slice []string, element string) int {
//	for index, elem := range slice {
//		if elem == element {
//			return index
//		}
//	}
//	return -1
//}
//
//// containsString returns true iff slice contains element
//func containsString(slice []string, element string) bool {
//	return !(posString(slice, element) == -1)
//}
