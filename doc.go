/*
.

Count and encrypt email addresses in files faster.

WHAT'S THE PURPOSE OF THIS SCRIPT?

- Count unique email addresses in a file. Ideal to work from a file produced by combining multiple csv and plain text files in multiple formats.

- Create lists of hashed (encrypted) email addresses that can be uploaded to Google Adwords. Tested up to 2.000.000 email addresses without issues.

- Beside email addresses, it can also count unique sha256 hashes in a file. It can be used to count unique email addresses when they are hashed.

GET HELP

./ecounter --help

.
*/
package main
