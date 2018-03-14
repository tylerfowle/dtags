package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ryanuber/columnize"
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
		if err := database.DeleteKey(strings.ToLower(args[0])); err != nil {
			fmt.Println("Failed.")
		} else {
			fmt.Printf("Successfully deleted tag [%s]\n", args[0])
		}
		break
	case "list":
		database.GetCurrentTags()
		break
	case "tags", "completion":
		printAllTags()
		break
	case "ls":
		printBoth()
		break
	default:
		args := os.Args[1:]
		printPath(args)
	}
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

	if err := database.AddKey(k, v); err != nil {
		fmt.Printf("Failed.\n")
	} else {
		fmt.Printf("Successfully added tag [%s] to path [%s]\n", k, v)
	}

}

func printAllTags() {
	for _, tag := range database.GetTags() {
		fmt.Println(tag)
	}
}

func printBoth() {
	var unformattedlist []string
	for tag, path := range database.All() {
		unformattedlist = append(unformattedlist, fmt.Sprintf("%s|%s\n", tag, path))
	}
	sort.Strings(unformattedlist)
	formattedList := columnize.SimpleFormat(unformattedlist)
	fmt.Println(formattedList)
}

func printPath(args []string) {
	cwd := database.GetValue(args[0])
	if cwd == "" {
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, cwd)
	os.Exit(1)
}

func confirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	switch response {
	case "y", "yes":
		fmt.Printf("Overwriting tag [%s]\n", os.Args[1])
		return true
	case "n", "no":
		fmt.Println("Overwrite cancelled.")
		return false
	default:
		fmt.Println("(yes/no) required:")
		return confirmation()
	}
}
