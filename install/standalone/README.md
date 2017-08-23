# Standalone setup

Below are instructions on how to set up the lexicon server for standalone use. For developer setup, see the `developer` folder.

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


## II. Installation

1. Download the install script: [install.sh](https://raw.githubusercontent.com/stts-se/pronlex/master/install/standalone/install.sh)

2. Install the pronlex server

     `$ sh install.sh <APPDIR>`

   Installs the pronlex server 


3. Import lexicon data (optional)

    `$ sh <APPDIR>/import.sh <APPDIR>`

   Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic.


## III. Start the lexicon server

The server should be started using this script

`$ sh <APPDIR>/start_server.sh <APPDIR>`

The startup script will run some init tests in a separate test server, before starting the standard server.

When the standard (non-testing) server is started, it always creates a demo database and lexicon, containing a few simple entries for demo and testing purposes. The server can thus be started and tested even if you haven't imported the lexicon data above.
