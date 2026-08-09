package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const emailRegex string = `([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)`

const shaRegex string = `[A-Fa-f0-9]{64}`

const dninieRegex string = `[A-z]?\d{7,8}[TRWAGMYFPDXBNJZSQVHLCKEtrwagmyfpdxbnjzsqvhlcke]`

const urlsRegex string = `https?://([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?`

var debug *bool

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	defer timeTrack(time.Now(), "main")

	fs := flag.NewFlagSet("ecounter", flag.ContinueOnError)
	help := fs.Bool("help", false, "Display help")
	inputFile := fs.String("input", "test.txt", "File or folder to do the operations")
	countIt := fs.String("count", "emails", "What to count")
	outputFile := fs.String("output", "output.txt", "File to output the results")
	debug = fs.Bool("debug", false, "Debug the script")
	encrypt := fs.Bool("encrypt", false, "Encrypts the emails as sha256")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}

	if *help {

		helpMe()
		return 0

	}

	if isDirectory(*inputFile) && outputInsideFolder(*inputFile, *outputFile) {
		fmt.Println("ERROR: The output file cannot be located inside the scanned folder or any of its subfolders:", *inputFile)
		return 1
	}

	fileToHandle := inputToString(*inputFile)

	var allMatches []string
	var allMatchesLC []string
	var allMatchesUC []string
	var uniques []string

	switch *countIt {
	case "emails":
		allMatches = searchInString(fileToHandle, emailRegex)
		allMatchesLC = arrayToLowercase(allMatches)
		uniques = uniquesInArray(allMatchesLC)
	case "sha256":
		allMatches = searchInString(fileToHandle, shaRegex)
		allMatchesLC = arrayToLowercase(allMatches)
		uniques = uniquesInArray(allMatchesLC)
	case "urls":
		allMatches = searchInString(fileToHandle, urlsRegex)
		uniques = uniquesInArray(allMatches)
	case "dnis":
		allMatches = searchInString(fileToHandle, dninieRegex)
		allMatchesUC = arrayToUpercase(allMatches)
		uniques = uniquesInArray(allMatchesUC)
	default:
		allMatches = searchInString(fileToHandle, emailRegex)
		allMatchesLC = arrayToLowercase(allMatches)
		uniques = uniquesInArray(allMatchesLC)
	}

	if *encrypt {
		uniques = arrayToSha256(uniques)
	}

	stringFinal := strings.Join(uniques, "\n")

	stringToFile(*outputFile, stringFinal)

	fmt.Println("\nWHAT HAPPENED?")
	fmt.Println("The parsed input : ", *inputFile)
	fmt.Println("Number of total", *countIt, "found in", *inputFile, ":", len(allMatches))
	fmt.Println("Number of unique", *countIt, "saved in the file", *outputFile, ":", len(uniques))
	fmt.Println("The results are hashed as sha256 ?", *encrypt)
	fmt.Printf("\n")

	return 0
}
