/*
Package dbapi contains code wrapped around SQL(ite3) and MariaDB.
It is used for inserting, updating and retrieving lexical entries from
a pronunciation lexicon database. A lexical entry is represented by
the lex.Entry struct, that mirrors entries of the entry database
table, along with associated tables such as transcription and lemma.
*/
package dbapi
