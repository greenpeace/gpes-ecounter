package main

import "fmt"

func helpMe() {

	textToPrint := `

* * * HELP * * *

Script to extract unique emails from a text file, including csv, txt or html. Optionally save it as sha 256.

Use the options as in this example:

./ecounter -input=rawfile.csv -output=uniques.txt -encripted=true

-help				Display this help
-input=rawfile.csv		Define the input file as rawfile.csv
-count=emails			What to count in the file. It can be "emails" or "sha256". By default it counts emails.
-output=uniques.txt		Define the output file as uniques.txt 
-encrypt=true			Encrypt the results as sha256, to upload to Google Adwords
-debug=true			Debug the script					

`
	fmt.Printf(textToPrint)

}
