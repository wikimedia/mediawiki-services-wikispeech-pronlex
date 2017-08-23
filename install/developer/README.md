# Setup for developers

Below are instructions on how to set up the lexicon server for development. For standalone setup, see the `standalone` folder.

## I. Preparation steps

1. Prerequisites

     If you're on Linux, you may need to install `gcc` and `build-essential` for the `sqlite3` go adapter to work properly:   
     `$ sudo apt-get install gcc`   
     `$ sudo apt-get install build-essential`

2. Set up `go`

     Download: https://golang.org/dl/ (1.8 or higher)   
     Installation instructions: https://golang.org/doc/install        
     Make sure the GOPATH variable is set: `$ echo $GOPATH` 

3. Install [Sqlite3](https://www.sqlite.org/)



## II. Download and install the pronlex library

1. Clone the source code

   `$ mkdir -p $GOPATH/src/github.com/stts-se`   
   `$ cd $GOPATH/src/github.com/stts-se`   
   `$ git clone https://github.com/stts-se/pronlex.git`


2. Download dependencies
   
   `$ cd $GOPATH/src/github.com/stts-se/pronlex`   
   `pronlex$ go get ./...`


3. Setup the pronlex server

   `$ sh setup.sh <APPDIR>`

   Installs files needed for the pronlex server


4. Import lexicon data (optional)

   `$ sh import.sh <APPDIR>`   

   Imports lexicon files for Swedish, Norwegian, US English and a small test file Arabic.


## III. Start the lexicon server

The server should be started using this script

`$ sh start_server.sh <APPDIR>`

The startup script will run some init tests in a separate test server, before starting the standard server.

When the standard (non-testing) server is started, it always creates a demo database and lexicon, containing a few simple entries for demo and testing purposes. The server can thus be started and tested even if you haven't imported the lexicon data above.

