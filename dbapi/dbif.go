package dbapi

import (
	"database/sql"

	"github.com/stts-se/pronlex/lex"
)

// DBIE is an interface that contains methods that make up the db api. This interface is used to create db connetions to different databases.
// NOTE: This interface is HUGE and may warrant refactoring. It was created to be able to add MariaDB support in addition to the original Sqlite (an afterthought).
type DBIF interface {

	// TODO Go through list below, maybe some function should be
	// removed, others be added

	associateLemma2Entry(db *sql.Tx, l lex.Lemma, e lex.Entry) error
	defineLexicon(db *sql.DB, l lexicon) (lexicon, error)
	deleteEntry(db *sql.DB, entryID int64, lexName string) (int64, error)
	deleteLexicon(db *sql.DB, lexName string) error
	entryCount(db *sql.DB, lexiconName string) (int64, error)
	getEntryFromID(db *sql.DB, id int64) (lex.Entry, error)
	getLexicon(db *sql.DB, name string) (lexicon, error)
	getLexiconMapTx(tx *sql.Tx) (map[string]bool, error)
	getLexiconTx(tx *sql.Tx, name string) (lexicon, error)
	insertEntries(db *sql.DB, l lexicon, es []lex.Entry) ([]int64, error)
	insertEntryComments(tx *sql.Tx, eID int64, eComments []lex.EntryComment) error
	insertEntryTagTx(tx *sql.Tx, entryID int64, tag string) error
	insertEntryValidations(tx *sql.Tx, e lex.Entry, eValis []lex.EntryValidation) error
	insertLemma(tx *sql.Tx, l lex.Lemma) (lex.Lemma, error)
	lexiconStats(db *sql.DB, lexName string) (LexStats, error)
	listAllEntryStatuses(db *sql.DB, lexiconName string) ([]string, error)
	listCommentLabels(db *sql.DB, lexiconName string) ([]string, error)
	listCurrentEntryStatuses(db *sql.DB, lexiconName string) ([]string, error)
	listCurrentEntryStatusesWithFreq(db *sql.DB, lexiconName string) (map[string]int, error)
	listCurrentEntryUsers(db *sql.DB, lexiconName string) ([]string, error)
	listCurrentEntryUsersWithFreq(db *sql.DB, lexiconName string) (map[string]int, error)
	listEntryStatuses(db *sql.DB, lexiconName string, onlyCurrent bool) ([]string, error)
	listEntryStatusesWithFreq(db *sql.DB, lexiconName string, onlyCurrent bool) (map[string]int, error)
	listEntryUsers(db *sql.DB, lexiconName string, onlyCurrent bool) ([]string, error)
	listEntryUsersWithFreq(db *sql.DB, lexiconName string, onlyCurrent bool) (map[string]int, error)
	firstTimePopulateDBCache(dbClusterLocation string) error
	listLexicons(db *sql.DB) ([]lexicon, error)
	locale(db *sql.DB, lexiconName string) (string, error)
	lookUp(db *sql.DB, lexNames []lex.LexName, q Query, out lex.EntryWriter) error
	lookUpIds(db *sql.DB, lexNames []lex.LexName, q Query) ([]int64, error)
	lookUpIdsTx(tx *sql.Tx, lexNames []lex.LexName, q Query) ([]int64, error)
	lookUpIntoMap(db *sql.DB, lexNames []lex.LexName, q Query) (map[string][]lex.Entry, error)
	lookUpIntoSlice(db *sql.DB, lexNames []lex.LexName, q Query) ([]lex.Entry, error)
	lookUpTx(tx *sql.Tx, lexNames []lex.LexName, q Query, out lex.EntryWriter) error
	moveNewEntries(db *sql.DB, fromLexicon, toLexicon, newSource, newStatus string) (MoveResult, error)
	moveNewEntriesTx(tx *sql.Tx, fromLexicon, toLexicon, newSource, newStatus string) (MoveResult, error)
	setOrGetLemma(tx *sql.Tx, strn string, reading string, paradigm string) (lex.Lemma, error)
	updateEntryComments(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	updateEntry(db *sql.DB, e lex.Entry) (res lex.Entry, updated bool, err error)
	updateEntryStatus(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (updated bool, err error)
	updateEntryTag(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	updateEntryTx(tx *sql.Tx, e lex.Entry) (updated bool, err error)
	updateEntryValidationForce(tx *sql.Tx, e lex.Entry) (bool, error)
	updateEntryValidation(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	updateLanguage(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	updateLemma(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (updated bool, err error)
	updateMorphology(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	updatePartOfSpeech(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	updatePreferred(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	updateTranscriptions(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (updated bool, err error)
	updateValidation(db *sql.DB, entries []lex.Entry) error
	updateValidationTx(tx *sql.Tx, entries []lex.Entry) error
	updateWordParts(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error)
	validateInputLexicons(tx *sql.Tx, lexNames []lex.LexName, q Query) error
	validationStats(db *sql.DB, lexName string) (ValStats, error)
	validationStatsTx(tx *sql.Tx, lexiconID int64) (ValStats, error)

	name() string
	engine() DBEngine
}

type DBEngine int

const (
	Sqlite DBEngine = iota

	MariaDB
)
