package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const emailRegex string = `([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)`

const shaRegex string = `[A-Fa-f0-9]{64}`

var debug *bool

func main() {

	defer timeTrack(time.Now(), "main")

	help := flag.Bool("help", false, "Display help")
	inputFile := flag.String("input", "test.txt", "File to do the operations")
	countIt := flag.String("count", "emails", "What to count")
	outputFile := flag.String("output", "output.txt", "File to output the results")
	debug = flag.Bool("debug", false, "Debug the script")
	encrypt := flag.Bool("encrypt", false, "Encrypts the emails as sha256")
	flag.Parse()

	if *help == true {

		helpMe()

	} else {

		fileToHandle := fileToString(*inputFile)

		var allMatches []string

		switch *countIt {
		case "emails":
			allMatches = searchInString(fileToHandle, emailRegex)
		case "sha256":
			allMatches = searchInString(fileToHandle, shaRegex)
		default:
			allMatches = searchInString(fileToHandle, emailRegex)
		}

		allMatchesLC := arrayToLowercase(allMatches)
		uniques := uniquesInArray(allMatchesLC)

		if *encrypt == true {
			uniques = arrayToSha256(uniques)
		}

		stringFinal := strings.Join(uniques, "\n")

		stringToFile(*outputFile, stringFinal)

		fmt.Println("\nWHAT HAPPENED?")
		fmt.Println("The parsed file : ", *inputFile)
		fmt.Println("Number of total", *countIt, "found in", *inputFile, ":", len(allMatches))
		fmt.Println("Number of unique", *countIt, "saved in the file", *outputFile, ":", len(uniques))
		fmt.Println("The results are hashed as sha256 ?", *encrypt)
		fmt.Printf("\n")

	}

}
