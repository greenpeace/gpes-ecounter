/*
.

Count and encrypt email addresses in files faster.

WHAT'S THE PURPOSE OF THIS SCRIPT?

- Count unique email addresses in a file. Ideal to work from a file produced by combining multiple csv and plain text files in multiple formats.

- Create lists of hashed (encrypted) email addresses that can be uploaded to Google Adwords. Tested up to 2.000.000 email addresses without issues.

- Beside email addresses, it can also count unique sha256 hashes in a file. It can be used to count unique email addresses when they are hashed.

GET HELP

	./ecounter --help

HOW TO USE IT

Scan emails in a file

Scan a file and create another file with just one email per line:

    ./ecounter -input=rawfile.csv -output=uniques.txt

Scan sha256 hashes in a file

Scan a file for sha256 hashes and create another file with just one (unique) sha256 hash per line:

    ./ecounter -input=rawfile.csv -count=sha256 -output=uniques.txt

This feature is useful to count unique emails when they are hashed (encrypted).

Scan urls in a file

Scan a file, like for example a sitemap.xml file, for urls:

    ./ecounter -input=sitemap.xml -count=urls -output=urls.csv

This feature is useful to create files with urls to use with check-my-pages.

Hash emails (encrypt)

    ./ecounter -input=rawfile.csv -output=uniques.txt -encrypt=true

Work from multiple files

To extract unique emails from multiple files you must concatenate the files in one plain .txt file first. Use the command line (cat in Linux or Mac and type in Windows) to join the csv or txt files in a combined.txt file:

    cat raw1.csv raw2.csv > combined.txt

Or to do it from all csv files in the all folder:

    cat all/*.csv > combined.txt

Please note that the files don't need to be in the same format. They just need to be csv or text files containing, among any information, email addresses. When the email addresses are all in one file, use it normally:

    ./ecounter -input=combined.txt -output=uniques.txt -encrypt=true

.
*/
package main
