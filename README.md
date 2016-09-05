# pronlex
pronlex is a pronunciation lexicon database behind a http API.

You need [Go](https://golang.org/) (1.7)

You need [Sqlite3](https://www.sqlite.org/)

----


Clone pronlex under src/github.com/stts-se/ in your [GOPATH](https://golang.org/doc/code.html#GOPATH) root.

```
cd pronlex/createEmptyDB/
go get
go install
cd ../importLexToDB
go install
cd ../lexserver
go install


Make sure GOPATH/bin is in $PATH
```

---
Create a pronlex.db and place it in the lexserver directory.

```
github.com/stts-se/pronlex/lexserver$ createEmptyDB pronlex.db
github.com/stts-se/pronlex/lexserver$ importLexToDB pronlex.db sv-se.nst swe030224NST.pron-ws.utf8.gz sv-se_ws-sampa sv-se_ws-sampa.csv
github.com/stts-se/pronlex/lexserver$ go run lexserver.go
```
(Or `CompileDaemon -command lexserver` if you use [Compile Daemon](https://github.com/githubnemo/CompileDaemon))


---
# Regexp db search performance

Regular expression search using a Go's regular expressions through the Sqlite3 driver is very slow. Either we should change databases, or find a better way to do regexp search in Sqlite3 from Go.


