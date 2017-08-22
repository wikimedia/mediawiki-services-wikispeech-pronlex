# Import lexicon files from the command line

## Create an empty database
    
      $ cd $GOPATH/src/github.com/stts-se/pronlex
      pronlex$ mkdir lexserver/db_files
      pronlex$ go run cmd/lexio/createEmptyDB/createEmptyDB.go lexserver/db_files/pronlex.db

## Import the files

You should be able to import the lexicon data below. Please note that it may take some time to load some of the larger lexicon files.

Before you start, `cd` to this folder: `$GOPATH/src/github.com/stts-se/pronlex`

#### Help and usage info

    pronlex$ go run cmd/lexio/import/import.go -help

#### Swedish

    pronlex$ go run cmd/lexio/import/import.go lexserver/db_files/pronlex.db sv-se.nst ~/gitrepos/lexdata/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa lexserver/symbol_set_file_area

#### Norwegian Bokmål

    pronlex$ go run cmd/lexio/import/import.go lexserver/db_files/pronlex.db nb-no.nst ~/gitrepos/lexdata/nb-no/nst/nor030224NST.pron-ws.utf8.gz nb-no_ws-sampa lexserver/symbol_set_file_area

#### US English

    pronlex$ go run cmd/lexio/import/import.go lexserver/db_files/pronlex.db en-us.cmu ~/gitrepos/lexdata/en-us/cmudict/cmudict-0.7b-ws.utf8 en-us_ws-sampa lexserver/symbol_set_file_area

#### Test Arabic

    pronlex$ go run cmd/lexio/import/import.go lexserver/db_files/pronlex.db ar-test ~/gitrepos/lexdata/ar/TEST/ar_TEST.pron-ws.utf8 ar_ws-sampa $GOPATH/src/github.com/stts-se/pronlex/lexserver/symbol_set_file_area

***
### Flags

Please note that flags have to be inserted _before_ the other arguments.

#### `replace` flag

Using the flag `-replace`, any pre-existing lexicon with the same name will be deleted before the specified file is loaded. 

#### `validate` flag

The instructions above will load lexicon files without validation. If you want to include validation, add the flag `-validate` to your command.

Example exec (Swedish):  

     pronlex$ go run cmd/lexio/import/import.go -validate lexserver/db_files/pronlex.db sv-se.nst ~/gitrepos/lexdata/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa lexserver/symbol_set_file_area

Example exec (US English):  

    pronlex$ go run cmd/lexio/import/import.go -validate lexserver/db_files/pronlex.db en-us.cmu ~/gitrepos/lexdata/en-us/cmudict/cmudict-0.7b-ws.utf8 en-us_ws-sampa lexserver/symbol_set_file_area


This will validate each entry according to the validation rules of it's associated project/language, and add the validation messages to the database.

#### `quiet` flag (in future releases)

This flag will stop the import command from printing continuous info on the amount of entries loaded (which may clog your log files). Committed to master as of 2017-08-18.


***

Prerequirements: https://github.com/stts-se/lexdata/wiki/Lexserver-setup-for-developers
