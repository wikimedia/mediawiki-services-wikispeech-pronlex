# pronlex
pronlex is a pronunciation lexicon database behind an http API.

[![GoDoc](https://godoc.org/github.com/stts-se/pronlex?status.svg)](https://godoc.org/github.com/stts-se/pronlex)

You need [Go](https://golang.org/) (1.7 or higher)  
You need [Sqlite3](https://www.sqlite.org/)

----
# Lexicon data and database

How to set up: https://github.com/stts-se/lexdata/wiki/Create-lexicon-database


---
# Regexp db search performance

Regular expression search using a Go's regular expressions through the Sqlite3 driver is very slow. Either we should change databases, or find a better way to do regexp search in Sqlite3 from Go.


