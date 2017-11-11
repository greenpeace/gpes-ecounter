package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// timeTrack is used to debug each function by measuring how long it takes to execute.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	if *debug == true {
		log.Printf("%s took %s", name, elapsed)
	}
}

// fileToString reads a file into a sting.
func fileToString(fileName string) string {
	defer timeTrack(time.Now(), "fileToString")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("ERROR: The file/path", fileName, "does not exist here")
		os.Exit(-1)
	}
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

// searchInString reads a string and returns all matches in the string.
func searchInString(total string, expression string) []string {
	defer timeTrack(time.Now(), "searchInString")
	r, err := regexp.Compile(expression)
	if err != nil {
		panic(err)
	}
	allMatches := r.FindAllString(total, -1)
	return allMatches
}

// arrayToLowercase Lowercases all the items in an array.
func arrayToLowercase(a []string) []string {
	defer timeTrack(time.Now(), "arrayToLowercase")
	var result []string
	for _, v := range a {
		lv := strings.ToLower(v)
		result = append(result, lv)
	}
	return result
}

//uniquesInArray Finds uniques in an array and returns.
func uniquesInArray(a []string) []string {
	defer timeTrack(time.Now(), "uniquesInArray")
	set := make(map[string]struct{})
	for _, v := range a {
		set[v] = struct{}{}
	}
	var uniques []string
	for k := range set {
		uniques = append(uniques, k)
	}
	return uniques
}

// stringToSha256 converts a string to sha 256.
func stringToSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	str := hex.EncodeToString(bs)
	return str
}

// arrayToSha256 Iterates a slice and converts it to sha256.
func arrayToSha256(a []string) []string {
	defer timeTrack(time.Now(), "arrayToSha256")
	var result []string
	var encodedSha string
	for _, v := range a {
		encodedSha = stringToSha256(v)
		result = append(result, encodedSha)
	}
	return result
}

// stringToFile writes a string to a file.
func stringToFile(fileName string, dat string) {
	defer timeTrack(time.Now(), "stringToFile")
	err := ioutil.WriteFile(fileName, []byte(dat), 0644)
	if err != nil {
		panic(err)
	}
}
