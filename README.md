# pronlex
pronlex is a pronunciation lexicon database with a server behind a simple HTTP API.

[![GoDoc](https://godoc.org/github.com/stts-se/pronlex?status.svg)](https://godoc.org/github.com/stts-se/pronlex)
[![Go Report Card](https://goreportcard.com/badge/github.com/stts-se/pronlex)](https://goreportcard.com/report/github.com/stts-se/pronlex) [![Build Status](https://travis-ci.org/stts-se/pronlex.svg?branch=master)](https://travis-ci.org/stts-se/pronlex)



## Lexicon server / Installation instructions

Utility scripts below (setup, import, start_server) require a working `bash` installation (preferably on a Linux system).

### I. Installation

1. Prerequisites

     If you're on Linux, you may need to install `gcc` and `build-essential` for the `sqlite3` go adapter to work properly:   
     `$ sudo apt-get install gcc build-essential`

2. Set up `go`

     Download: https://golang.org/dl/ (1.13 or higher)   
     Installation instructions: https://golang.org/doc/install             


3. Install [Sqlite3](https://www.sqlite.org/)

     On Linux systems with `apt`, run `sudo apt install sqlite3`


4. Clone the source code

   `$ git clone https://github.com/stts-se/pronlex.git`  
   `$ cd pronlex`   
   
5. Test (optional)

   `pronlex$ go test ./...`


### II. Quick start: Create a lexicon database file and look up a word

1) Download an SQL lexicon dump file. In the following example, we use a Swedish lexicon: `https://github.com/stts-se/lexdata/blob/master/sv-se/nst/swe030224NST.pron-ws.utf8.sql.gz`

2) Pre-compile binaries (for faster execution times)

    `pronlex$ go build ./...`

2) Create a database file (this takes a while):

    `pronlex$ importSql swe030224NST.pron-ws.utf8.sql.gz swe_lex.db`
       
3) Test looking up a word:
       
   `pronlex$ lexlookup swe_lex.db åsna`


### III: Server setup

1. Setup the pronlex server

   `pronlex$ cd install`   
   `install$ bash setup.sh <APPDIR>`   
   Example:
   `install$ bash setup.sh ~/wikispeech`

   Sets up the pronlex server and a set of test data in the folder specified by `<APPDIR>`.


2. Import lexicon data (optional)

   `install$ bash import.sh <LEXDATA-GIT> <APPDIR>`   
   Example:
   `install$ bash import.sh ~/git_repos/lexdata ~/wikispeech` 

   Imports lexicon databases (sql dumps) for Swedish, Norwegian, US English, and a small set of test data for Arabic from the [lexdata repository](https://github.com/stts-se/lexdata).
If the `<LEXDATA-GIT>` folder exists on disk, lexicon resources will be read from this folder. If it doesn't exist, the lexicon data will be downloaded from github.

   If you want to import other lexicon data, or just a subset of the data above, you can use one of the following methods:
   
   * Import lexicon files using the lexserver web API (http://localhost:8787/admin/lex_import_page if you have a running lexicon server on localhost)
   * Import lexicon files from the command line using this import script: https://github.com/stts-se/pronlex/tree/master/cmd/lexio/importLex.
   * Import database sql dumps files from the command line using this import script: https://github.com/stts-se/pronlex/tree/master/cmd/lexio/importSql.


You can create your own lexicon files, or you can use existing data in the [lexdata repository](https://github.com/stts-se/lexdata). The lexicon file format is described here: https://godoc.org/github.com/stts-se/pronlex/line.


### IV. Start the lexicon server

The server is started using this script

`install$ bash start_server.sh -a <APPDIR>`

The startup script will run some init tests in a separate test server, before starting the standard server.

When the standard (non-testing) server is started, it always creates a demo database and lexicon, containing a few simple entries for demo and testing purposes. The server can thus be started and tested even if you haven't imported the lexicon data above.

To specify port, run:   
`install$ bash start_server.sh -a <APPDIR> -p <PORT>`


For a complete set of options, run:  
`install$ bash start_server.sh -h`


<!--

## For developers

If you are developing for Wikispeech, and need to make changes to this repository, make sure you run a test build using `build_and_test.sh` before you make a pull request. Don't run more than one instance of this script at once, and make sure no pronlex server is already running on the default port.

Wikimedia's installation instructions for Wikispeech: https://www.mediawiki.org/wiki/Extension:Wikispeech

-->


