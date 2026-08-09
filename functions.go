package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// timeTrack is used to debug each function by measuring how long it takes to execute.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	if *debug {
		log.Printf("%s took %s", name, elapsed)
	}
}

// fileToString reads a file into a sting.
func fileToString(fileName string) string {
	defer timeTrack(time.Now(), "fileToString")
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("ERROR: The file/path", fileName, "does not exist here")
		os.Exit(1)
	}
	dat, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(dat)
}

// isDirectory returns true if the given path is a directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// inputToString reads all the content from a file or, if the input is a folder, from all the supported files inside the folder and its subfolders.
func inputToString(inputPath string) string {
	defer timeTrack(time.Now(), "inputToString")
	if isDirectory(inputPath) {
		return folderToString(inputPath)
	}
	return fileToString(inputPath)
}

// supportedExtension returns true if the file has one of the supported extensions.
func supportedExtension(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".txt", ".csv", ".tsv", ".xml", ".html", ".json":
		return true
	}
	return false
}

// folderToString reads the content of all the supported files inside a folder and its subfolders.
func folderToString(folderPath string) string {
	var sb strings.Builder
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !supportedExtension(path) {
			return nil
		}
		dat, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sb.Write(dat)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return sb.String()
}

// outputInsideFolder returns true if the output file would be written inside the scanned folder or any of its subfolders.
func outputInsideFolder(folderPath, outputPath string) bool {
	folderAbs, err := filepath.Abs(folderPath)
	if err != nil {
		return false
	}
	outputAbs, err := filepath.Abs(outputPath)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(folderAbs, filepath.Dir(outputAbs))
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
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

// arrayToUpercase Upercases all the items in an array.
func arrayToUpercase(a []string) []string {
	defer timeTrack(time.Now(), "arrayToUpercase")
	var result []string
	for _, v := range a {
		lv := strings.ToUpper(v)
		result = append(result, lv)
	}
	return result
}

// uniquesInArray Finds uniques in an array and returns.
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
	err := os.WriteFile(fileName, []byte(dat), 0644)
	if err != nil {
		panic(err)
	}
}
