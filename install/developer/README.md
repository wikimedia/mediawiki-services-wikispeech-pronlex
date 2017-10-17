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
     
     Add your `GOPATH/bin` to your `PATH` environment variable.    
     Your `GOPATH` can be found using the following command:    
     `$ go env GOPATH`

     On Linux-like systems, this is typically the command you could add to your `.profile` file or similar:    
     `export PATH=$PATH:$(go env GOPATH)/bin`

3. Install [Sqlite3](https://www.sqlite.org/)



## II. Installation

1. Clone the source code

   `$ mkdir -p $(go env GOPATH)/src/github.com/stts-se`   
   `$ cd $(go env GOPATH)/src/github.com/stts-se`   
   `stts-se$ git clone https://github.com/stts-se/pronlex.git`


2. Download dependencies
   
   `$ cd $(go env GOPATH)/src/github.com/stts-se/pronlex`   
   `pronlex$ go get ./...`


3. Setup the pronlex server

   `pronlex$ cd install/developer`   
   `developer$ sh setup.sh <APPDIR>`

   Sets up the pronlex server and a set of test data in the folder specified by `<APPDIR>`.


4. Import lexicon data (optional)

   `developer$ sh import.sh <LEXDATA-GIT> <APPDIR>`   

   Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic from the lexdata repository.
If the `<LEXDATA-GIT>` folder doesn't exist, it will be downloaded from github.


## III. Start the lexicon server

The server is started using this script

`developer$ sh start_server.sh -a <APPDIR>`

The startup script will run some init tests in a separate test server, before starting the standard server.

When the standard (non-testing) server is started, it always creates a demo database and lexicon, containing a few simple entries for demo and testing purposes. The server can thus be started and tested even if you haven't imported the lexicon data above.

To specify port, run:   
`$ sh start_server.sh -a <APPDIR> -p <PORT>`


For a complete set of options, run:  
`$ sh start_server.sh -h`

