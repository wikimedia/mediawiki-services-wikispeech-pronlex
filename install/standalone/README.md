# Standalone setup

Below are instructions on how to set up the lexicon server for standalone use. For developer setup, see the `developer` folder.

## I. Preparation steps

1. Prerequisites

     If you're on Linux, you may need to install `gcc` and `build-essential` for the `sqlite3` go adapter to work properly:   
     `$ sudo apt-get install gcc`   
     `$ sudo apt-get install build-essential`

2. Set up `go` and `sqlite3`

     1. Install `go`   
        Download: https://golang.org/dl/ (1.8 or higher)   
        Installation instructions: https://golang.org/doc/install
        
     2. Install [Sqlite3](https://www.sqlite.org/)


## II. Installation scripts

1. Install the pronlex server

     `$ sh install.sh <APPDIR>` | [download script](https://raw.githubusercontent.com/stts-se/pronlex/master/install/standalone/install.sh)

   Installs the pronlex server 


2. Import lexicon data (optional)

    `$ sh import.sh <APPDIR>` | [download script](https://raw.githubusercontent.com/stts-se/pronlex/master/install/standalone/import.sh)

   Imports lexicon files for Swedish, Norwegian, US English and a small test file Arabic.


## III. Start the lexicon server

The server should be started using this script

`$ sh start_server.sh <APPDIR>` | [download script](https://raw.githubusercontent.com/stts-se/pronlex/master/install/standalone/start_server.sh)

The startup script will run some init tests in a separate test server, before starting the standard server.

When the standard (non-testing) server is started, it always creates a demo database and lexicon, containing a few simple entries for demo and testing purposes. The server can thus be started and tested even if you haven't imported the lexicon data above.
