package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

const emailRegex string = `([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)`

var debug *bool

func main() {

	defer timeTrack(time.Now(), "main")

	help := flag.Bool("help", false, "Display help")
	inputFile := flag.String("input", "test.txt", "File to do the operations")
	outputFile := flag.String("output", "output.txt", "File to output the results")
	debug = flag.Bool("debug", false, "Display help")
	encrypted := flag.Bool("encrypted", false, "Encrypts the emails as sha256")
	flag.Parse()

	if *help == true {

		helpMe()

	} else {

		fileToHandle := fileToString(*inputFile)

		allMatches := searchInString(fileToHandle, emailRegex)

		allMatchesLC := arrayToLowercase(allMatches)
		uniques := uniquesInArray(allMatchesLC)

		if *encrypted == true {
			uniques = arrayToSha256(uniques)
		}

		stringFinal := strings.Join(uniques, "\n")

		stringToFile(*outputFile, stringFinal)

		fmt.Println("\nWHAT HAPPENED?")
		fmt.Println("The parsed file is: ", *inputFile)
		fmt.Println("The results are in the file: ", *outputFile)
		fmt.Println("Number of non unique emails found in", *inputFile, ":", len(allMatches))
		fmt.Println("Number of unique emails in", *outputFile, ":", len(uniques))
		fmt.Println("Are the results encrypted with sha256?", *encrypted)
		fmt.Printf("\n")

	}

}
