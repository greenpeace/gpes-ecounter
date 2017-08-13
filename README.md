# Ecounter

**Count and encrypt email addresses in files faster.**

## What's the purpose of this script?

* **Count unique email addresses** in a file. Ideal to work from a file produced by combining multiple csv and plain text files in multiple formats.
* Create lists of hashed (encrypted) email addresses that can be [uploaded to **Google Adwords**](https://support.google.com/adwords/answer/6276125?hl=en). Tested up to 2.000.000 email addresses without issues.
* Beside email addresses, it can also **count unique sha256** hashes in a file. It can be used to count unique email addresses when they are hashed.

## Why use it?

* **It's fast** - It parses a 75MB file, counting, selecting and hashing (encrypting) more than 1.000.000 unique email addresses in about 10 seconds in a good laptop computer.
* **It's easy to use** - It parses any csv, html, sql or plain text file. It does not require a specific format and can work from a file combining different formats.
* **It's free** -  No cost and the freedom to study and modify to fit.

## How to use it

#### Get help

**ecounter** is a command line script. To get help type in the command line:

```
./ecounter --help
```

#### Scan emails in a file

Scan a file and create another file with just one email per line:

```
./ecounter -input=rawfile.csv -output=uniques.txt
```

#### Scan sha256 hashes in a file

Scan a file for sha256 hashes and create another file with just one (unique) sha256 hash per line:

```
./ecounter -input=rawfile.csv -count=sha256 -output=uniques.txt
```

This feature is useful to count unique emails when they are hashed (encrypted).

#### Hash emails (encrypt)

```
./ecounter -input=rawfile.csv -output=uniques.txt -encrypt=true
```

#### Counting

The script produces a report like this:

```
WHAT HAPPENED?
The parsed file :  rawfile.csv
Number of total emails found in rawfile.csv : 1420463
Number of unique emails saved in the file uniques.txt : 1329624
The results are hashed as sha256 ? true
```

#### Work from multiple files 

To extract unique emails from multiple files you must concatenate the files in one plain .txt file first. Use the command line (`cat` in Linux or Mac and `type` in Windows) to join the csv or txt files in a **combined.txt** file:

```
cat raw1.csv raw2.csv > combined.txt 
```

Or to do it from all **csv** files in the **all** folder:

```
cat all/*.csv > combined.txt
```

Please note that the files don't need to be in the same format. They just need to be csv or text files containing, among any information, email addresses. When the email addresses are all in one file, use it normally:

```
./ecounter -input=combined.txt -output=uniques.txt -encrypt=true
```

## Install

Currently this script is provided as source code in **[Go](https://golang.org/dl/)**. You can easily download it and compile. Still if some people find it useful I can compile executables and add the links here.

## Problems?

* If **the total number of emails doesn't match** it's because this script uses a regular expression to find emails in a file. Sometimes certain email addresses that are accepted by other softwares don't match this regular expression. To modify the regular expression change the constant `emailRegex` in the script. Try different regular expressions if you want, but it's likely that the unaccounted emails are not valid, so why should you count them?

