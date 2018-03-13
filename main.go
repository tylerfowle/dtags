package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tylerfowle/dtags/db"
)

var database *db.Database

func main() {
	var err error

	database, err = db.Init()
	if err != nil {
		panic(err)
	}

	defer database.Instance.Close()

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "add":
		addNewTag(args)
		break
	case "del":
		database.DeleteKey(strings.ToLower(args[0]))
		break
	case "list":
		printAllTags()
	}

	//case "list":
	//	listTags(info)
	//case "ls":
	//	listAll(info)
	//case "completion":
	//	tagCompletion(info)
	//default:
	//	info.key = []byte(info.subcommand)
	//	info.args = []string(os.Args[1:])
	//	shell(info)
	//}

}

func addNewTag(args []string) {
	k := strings.ToLower(args[0])
	v := database.CurrentDirectory

	if len(args[0:]) > 1 {
		v = args[1]
	}

	if database.Exists(k) {
		fmt.Printf("Overwrite existing tag? [%s] (y/n)", k)
		if confirmation() == false {
			return
		}
	}

	database.AddKey(k, v)
}

func printAllTags() {
	for _, tag := range database.GetTags() {
		fmt.Println(tag)
	}
}

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
func confirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	y := []string{"y", "Y", "yes", "Yes", "YES"}
	n := []string{"n", "N", "no", "No", "NO"}

	response = strings.TrimSpace(response)
	response = strings.ToLower(response)

	if containsString(y, response) {
		return true
	} else if containsString(n, response) {
		return false
	} else {
		fmt.Println("yes or no required:")
		return confirmation()
	}
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

//
// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
