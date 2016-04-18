# pronlex
pronlex is a sketch of/place holder for a pronunciation lexicon database behind a simple http API. It is NOT ready for proper use.

You need [Go](https://golang.org/)

You need [Sqlite3](https://www.sqlite.org/)

----


Clone pronlex under src/github.com/stts-se/ in your [GOPATH](https://golang.org/doc/code.html#GOPATH) root.

```
cd pronlex/createEmptyDB/
go get
go install
cd ../addNSTLexToDB
go install
cd ../lexserver
go install


Make sure GOPATH/bin is in $PATH
```

---
Create a pronlex.db and place it in the lexserver directory.

```
github.com/stts-se/pronlex/lexserver$ createEmptyDB pronlex.db
github.com/stts-se/pronlex/lexserver$ addNSTLexToDB sv.se.nst pronlex.db [PATH TO NST LEX]/swe030224NST.pron_utf8.txt 
github.com/stts-se/pronlex/lexserver$ go run lexserver.go
```
(Or `CompileDaemon -command lexserver` if you use [Compile Daemon](https://github.com/githubnemo/CompileDaemon))


---
# Regexp db search performance

Regular expression search using a Go's regular expressions through the Sqlite3 driver is very slow. Either we should change databases, or find a better way to do regexp search in Sqlite3 from Go.


