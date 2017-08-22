# Setup for developers

WORK IN PROGRESS. THESE INSTRUCTIONS ARE BEING UPDATED RIGHT NOW. PLEASE DO NOT USE.


Below are instructions on how to set up the lexicon server for development. For standalone setup, see STANDALONE.md.

## I. Preparation steps

1. Prerequisites

     If you're on Linux, you may need to install `gcc` and `build-essential` for the `sqlite3` go adapter to work properly (the adapter is installed in _Download dependencies_ below):   
     `$ sudo apt-get install gcc`   
     `$ sudo apt-get install build-essential`

2. Set up `go` and sqlite3

     1. Install `go` following the instructions here: https://golang.org/dl/ (1.8 or higher)

     2. Set your `$GOPATH` (we suggest `~/go`). Make sure the go binaries are in your `$PATH`.  
        If you're on Linux, add `export PATH=$PATH:/usr/local/go/bin:GOPATH/bin` to your `.bashrc` file.  
        (if you installed `go` here: `/usr/local/go`)

     3. Install [Sqlite3](https://www.sqlite.org/)

## II. Installation scripts (work in progress)

1. Install the pronlex server

   1. Clone the source code

    `$ mkdir -p $GOPATH/src/github.com/stts-se`  
    `$ cd $GOPATH/src/github.com/stts-se`  
    `stts-se$ git clone https://github.com/stts-se/pronlex.git`  

    2. Download dependencies
    
    `$ cd $GOPATH/src/github.com/stts-se/pronlex`   

    `pronlex$ go get ./...`   
      &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  or, if you want to know what's going on:    
    `pronlex$ go get -v ./...`

  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  Please note that the verbosity flag may give you a few confusing warnings, but you will at least see what packages are being processed.

     3. Clone the lexdata repository:
    
     `$ mkdir -p ~/gitrepos`  
     `$ cd ~/gitrepos`  
     `gitrepos$ git clone https://github.com/stts-se/lexdata.git`


     4. Prepare symbol sets and symbol set mappers/converters
    
     `$ cd $GOPATH/src/github.com/stts-se/pronlex/lexserver`
     `lexserver$ mkdir symbol_sets`  
     `lexserver$ cp ~/gitrepos/lexdata/*/*/*.sym symbol_sets`   
     `lexserver$ cp ~/gitrepos/lexdata/mappers.txt symbol_sets`  
     `lexserver$ cp ~/gitrepos/lexdata/converters/*.cnv symbol_sets`  

---
** UPCOMING: SCRIPT TO INSTALL pronlex INSTEAD OF COMMANDS **

`$ sh install.sh <APPDIR>`

   Installs the pronlex server and a small demo db for testing
---


2. Import lexicon files (optional)

`$ sh import.sh <APPDIR>`   

Imports lexicon files for Swedish, Norwegian, US English and a small test file Arabic.


## III. Start the lexicon server
The server should be started using this script

`sh start_server.sh`

The startup script will run some init tests in a separate test server, before starting the standard server.

When the standard (non-testing) server is started, it always creates a demo database and lexicon, containing a few simple entries for demo and testing purposes. The server can thus be started and tested even if you haven't imported the lexicon data provided on this site.

