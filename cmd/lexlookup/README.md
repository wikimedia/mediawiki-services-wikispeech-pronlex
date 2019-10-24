# lexlookup
Command line tool for look up in a [pronlex](https://github.com/stts-se/pronlex) Sqlite3 DB file.


    pronlex$ cd cmd/lexlookup/	

    go build
    ./lexlookup <PRONLEX SQLITE3 DB FILE> (WORDS | STDIN)

or (slower start up):

    go run main.go <PRONLEX SQLITE3 DB FILE> (WORDS | STDIN)




The `-missing` flag prints out words not found in the lexicon db, for example like this:

     cat words.txt | ./lexlookup pronlex.db -missing



There is also a `delete` flag for deleting entries from the db.