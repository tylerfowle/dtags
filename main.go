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

// Entrypoint for dtags application. Initializes database instance and
// processed passed command to perform the expected action.
func main() {
	var err error
	if database, err = db.Init(); err != nil {
		panic(err)
	}
	defer database.Instance.Close()

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		addNewTag(os.Args[2:])
		break
	case "del":
		database.DeleteKey(strings.ToLower(os.Args[2:][0]))
		break
	case "list", "completion":
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

// Show the help message to a user when not enough arguments are passed
// to the command.
func showHelp() {
	fmt.Println("usage: dtags command [params]")
	fmt.Println()
	fmt.Println("       add <tag> [path] - add a new tag at path; defaults to current directory")
	fmt.Println("       del <tag> - delete the provided tag from storage")
	fmt.Println("       list - list all tags currently stored")
	fmt.Println("       ls - list all tags currently stored and associated paths")
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
		fmt.Printf("tag not found\n")
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
	y := []string{"y", "yes"}
	n := []string{"n", "no"}

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

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
