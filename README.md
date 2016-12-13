# pronlex
pronlex is a pronunciation lexicon database behind a http API.

[![GoDoc](https://godoc.org/github.com/stts-se/pronlex?status.svg)](https://godoc.org/github.com/stts-se/pronlex)

You need [Go](https://golang.org/) (1.7)

You need [Sqlite3](https://www.sqlite.org/)

----


Obtain lexicon data files from the [lexdata](https://github.com/stts-se/lexdata) repository. Here you will also find symbol sets for each lexicon.

Clone pronlex under src/github.com/stts-se/ in your [GOPATH](https://golang.org/doc/code.html#GOPATH) root.

Make sure GOPATH/bin is in $PATH

---
Create a pronlex.db and place it in the lexserver directory.

```
github.com/stts-se/pronlex$ go run dbapi/createEmptyDB/createEmptyDB.go pronlex.db
github.com/stts-se/pronlex$ go run lexio/import/import.go pronlex.db sv-se.nst [LEX FILE FOLDER]/sv-se/nst/swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa [SYMBOLSET FOLDER]
# [SYMBOLSET FOLDER] is used for validation upon import. skip if you don't want validation (it takes some extra time).

github.com/stts-se/pronlex$ mv pronlex.db lexserver/

github.com/stts-se/pronlex$ cd lexserver

github.com/stts-se/pronlex/lexserver$ mkdir symbol_set_file_area
github.com/stts-se/pronlex/lexserver$ cp [LEXDATA]/*/*/*.tab symbol_set_file_area

github.com/stts-se/pronlex/lexserver$ go run *.go
```


---
# Regexp db search performance

Regular expression search using a Go's regular expressions through the Sqlite3 driver is very slow. Either we should change databases, or find a better way to do regexp search in Sqlite3 from Go.


