## pronlex
pronlex is a pronunciation lexicon database behind an http API.

[![GoDoc](https://godoc.org/github.com/stts-se/pronlex?status.svg)](https://godoc.org/github.com/stts-se/pronlex)


## Docker installation

TODO: add instructions here


## Lexicon server: Lexicon data and database

* **Standalone setup**    
File: standalone/README.md   
URL: https://github.com/stts-se/pronlex/blob/master/install/standalone/README.md

* **Setup for developers**    
File: developer/README.md   
URL: https://github.com/stts-se/pronlex/blob/master/install/developer/README.md



## Regexp db search performance

Regular expression search using a Go's regular expressions through the Sqlite3 driver is very slow. Either we should change databases, or find a better way to do regexp search in Sqlite3 from Go.


