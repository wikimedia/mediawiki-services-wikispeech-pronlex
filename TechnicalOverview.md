## Tehcnical overview of `pronlex` and the lexicon server

The basic function of `pronlex` is to store and retrieve lexical _entries_. An entry consist of a word form, along with a phonetic transcription, a status, a database and lexicon name, and possibly additional values.

A code version of an entry is defined in [lex.Entry](https://github.com/stts-se/pronlex/blob/master/lex/entry.go). Documentation is is available [here](https://godoc.org/github.com/stts-se/pronlex/lex).


An entry can be converted to and from JSON.


### Architecture

The `pronlex` package consists of a lexicon database and a lexicon server, plus some additional helper modules. The server is written in [`Go`](https://golang.org).

The lexicon server has an HTTP Rest API. More information about the HTTP API can be found below.

The lexicon database stores entries in a relational database, [Sqlite3](https://sqlite.org/index.html). The SQL schema --- the definition of the database structure --- is a string constant found in the file [schema.go](https://github.com/stts-se/pronlex/blob/master/dbapi/schema.go).



### HTTP API

There is an HTTP server for the pronlex database. A documentation of the HTTP API can be accessed once the server is started (default address: http://localhost:8787).

The most important API URLs can be found in the list below. For more information, and a complete list of API calls, please see the full documentation using local running lexicon server.

* /lexicon/list
* /lexicon/lookup
* /lexicon/entries_exist
* /lexicon/info/{lexicon_name}
* /lexicon/stats/{lexicon_name}
* /lexicon/updateentry
* /lexicon/addentry
* /lexicon/delete_entry/{lexicon_name}/{entry_id}
* /admin/list_dbs
* /admin/create_db/{db_name}
* /admin/define_lex/{lexicon_name}/{locale}/{symbolset_name}
* /admin/deletelexicon/{lexicon_name}
* /admin/superdeletelexicon/{lexicon_name}




### Database API


The database can be called using a set of functions defined in the database manager, [dbapi.DBManager](https://github.com/stts-se/pronlex/blob/master/dbapi/db_manager.go).

Internally, the database interaction is performed using functions defined in  [dbapi.go](https://github.com/stts-se/pronlex/blob/master/dbapi/dbapi.go).


The database can be queried through the `DBManager` using a query struct, [dbapi.DBMQuery](https://github.com/stts-se/pronlex/blob/master/dbapi/db_manager.go)


The DBMQuery contains the reference to a lexicon and the actual [dbapi.Query](https://godoc.org/github.com/stts-se/pronlex/dbapi#Query).

Such a query struct can be converted to and from JSON.


### Database queries

A query from the dbapi is converted to a SQL query string. This happens in [sql_gen.go](https://github.com/stts-se/pronlex/blob/master/dbapi/sql_gen.go).

The query string is then used to retrieve entries using functions in `dbapi`. 


### Helper commands

There are stand-alone commands for managing the lexicon database. These are located in the `cmd` folder.

* createEmptyDB - create an empty lexicon database (sqlite) file
* createEmptyLexicon - create an empty lexicon in a lexicon database
* exportLex - export a lexicon from a database file to a text file
* importLex - import a lexicon (text) file to a database
* importSql - import an lexicon sql dump into a database file
* lexlookup - command line tool for lexicon search/lookup
* validate_lex_file - command line tool for validating a lexicon (text) file


### Sqlite commands

 * Create an sql dump from a database:
`sqlite3 <dbFile> .dump | gzip -c > <sqlDumpFile>`

 * Import an sql dump to a database:
`gunzip -c <sqlDumpFile> | sqlite3 <dbFile>`


### List of additional modules and features

* symbolset - phonetic symbol sets are defined per lexicon, and can be used for validation
* symbolset mapper - component to convert between different phonetic symbol sets in the same language
* converter - component to convert transcription between different languages
* validation - validation components can be created to validate lexicon entries for various issues
  - phonetic symbols
  - transcription format
  - phonotactic rules
  - syllable boundaries
  - sanity checks
  - etc
* lexicon data conversion (for more information, see the [lexdata](https://github.com/stts-se/lexdata)) repository
  - CMU2WS - CMU US English
  - csCzPhword2WS - Czech dictionary
  - nbNoNST2WS - NST Norwegian bokmål
  - svSeNST2WS - NST Swedish
* admin - various admin tools apart from those listed above
