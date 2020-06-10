package dbapi

import (
	"strconv"
	"strings"

	"github.com/stts-se/pronlex/lex"
)

// ***** helpers -->

// trm trims spaces off a string
func trm(s string) string { return strings.TrimSpace(s) }

// nQs returns a string of n comma separated '?':s inside a pair of parens, such as "(?,?)"
func nQs(n int) string {
	var res string
	if n < 1 {
		return res
	}

	res += "(" + strings.TrimSuffix(strings.Repeat("?,", n), ",") + ")"
	return res
}

// convI convert a slice of ints into a slice of interfac
func convI(is []int64) []interface{} {
	res := make([]interface{}, len(is))
	for i, v := range is {
		res[i] = v
	}
	return res
}

// convS converts a slice of strings into a slice of interface
func convS(s []string) []interface{} {
	res := make([]interface{}, len(s))
	for i, v := range s {
		res[i] = v
	}
	return res
}

// <-- helpers *****

// tables returns at least 'entry' since it makes no sense to return the empty string,
// since the return value of this function is to be used after 'select
// Entry.id from'
// func tables(lexNames []lex.LexName, q Query) string {
// 	var res []string
// 	if len(lexNames) > 0 {
// 		res = append(res, "lexicon")
// 	}
// 	//if len(q.Words) > 0 || q.WordLike != "" || q.PartOfSpeechLike != "" {
// 	res = append(res, "Entry")
// 	//}

// 	if q.TranscriptionLike != "" {
// 		res = append(res, "Transcription")
// 	}
// 	if len(q.Lemmas) > 0 || q.LemmaLike != "" || q.ReadingLike != "" || q.ParadigmLike != "" {
// 		res = append(res, "lemma, Lemma2Entry")
// 	}

// 	return strings.Join(res, ", ")
// }

// lexicons returns a piece of sql matching the Querys list of
// lexicons and a slice of the db ids of the lexicons listed in a Query
func lexicons(lexNames []lex.LexName) (string, []interface{}) {
	var res string
	var resv []interface{}
	if len(lexNames) == 0 {
		return res, resv
	}

	res += "Entry.lexiconid = Lexicon.id AND Lexicon.name in " + nQs(len(lexNames))

	lNames := make([]interface{}, len(lexNames))
	for i, l := range lexNames {
		lNames[i] = string(l)
	}

	resv = append(resv, lNames...)
	return res, resv
}

func words(lexNames []lex.LexName, q Query) (string, []interface{}) {
	var reses []string
	var resv []interface{}

	// If none of the following values are set, there are no
	// references to entries in the query. This should not make
	// sense, since we are building a query to look for entries,
	// but such a reference can be add by searching for other
	// things depending on entry (such as transcription)

	// You must study the Query struct to understand this

	//fmt.Printf("sql_gen QUERY : %#v\n", q)

	if len(q.Words) == 0 && len(q.WordParts) == 0 && trm(q.WordLike) == "" && trm(q.WordPartsLike) == "" && trm(q.WordPartsRegexp) == "" && trm(q.WordRegexp) == "" && trm(q.PartOfSpeechLike) == "" && trm(q.PartOfSpeechRegexp) == "" && trm(q.LanguageLike) == "" && trm(q.MorphologyLike) == "" && len(q.EntryIDs) == 0 {
		return "", resv
	} //else {
	if len(q.Words) > 0 {
		reses = append(reses, "Entry.strn in "+nQs(len(q.Words)))
		resv = append(resv, convS(ToLower(q.Words))...)
	}
	if len(q.WordParts) > 0 {
		reses = append(reses, "Entry.wordParts in "+nQs(len(q.WordParts)))
		resv = append(resv, convS(ToLower(q.WordParts))...)
	}
	if trm(q.WordPartsLike) != "" {
		reses = append(reses, "Entry.wordParts like ?")
		resv = append(resv, q.WordPartsLike)
	}
	if trm(q.WordPartsRegexp) != "" {
		reses = append(reses, "Entry.wordParts REGEXP ?")
		resv = append(resv, q.WordPartsRegexp)
	}
	if len(q.EntryIDs) > 0 {
		reses = append(reses, "Entry.id in "+nQs(len(q.EntryIDs)))
		resv = append(resv, convI(q.EntryIDs)...)
	}
	if trm(q.WordLike) != "" {
		reses = append(reses, "Entry.strn like ?")
		resv = append(resv, q.WordLike)
	}
	if trm(q.WordRegexp) != "" {
		reses = append(reses, "Entry.strn REGEXP ?")
		resv = append(resv, q.WordRegexp)
	}

	if trm(q.LanguageLike) != "" {
		reses = append(reses, "Entry.language like ?")
		resv = append(resv, q.LanguageLike)
	}
	if trm(q.MorphologyLike) != "" {
		reses = append(reses, "Entry.morphology like ?")
		resv = append(resv, q.MorphologyLike)
	}
	if trm(q.PartOfSpeechLike) != "" {
		reses = append(reses, "Entry.partOfSpeech like ?")
		resv = append(resv, q.PartOfSpeechLike)
	}
	if trm(q.PartOfSpeechRegexp) != "" {
		reses = append(reses, "Entry.partOfSpeech REGEXP ?")
		resv = append(resv, q.PartOfSpeechRegexp)
	}

	//}

	res := strings.Join(reses, " AND ")

	if len(lexNames) != 0 {
		res += " and Entry.lexiconId = Lexicon.id"
	}

	return res, resv
}

func lemmas(q Query) (string, []interface{}) {
	var reses []string
	var resv []interface{}

	if len(q.Lemmas) == 0 && trm(q.LemmaLike) == "" && trm(q.LemmaRegexp) == "" &&
		trm(q.ReadingLike) == "" && trm(q.ReadingRegexp) == "" &&
		trm(q.ParadigmLike) == "" && trm(q.ParadigmRegexp) == "" {
		return "", resv
	}
	if len(q.Lemmas) > 0 {
		reses = append(reses, "Lemma.strn in "+nQs(len(q.Lemmas)))
		resv = append(resv, convS(q.Lemmas)...)
	}
	if trm(q.LemmaLike) != "" {
		reses = append(reses, "Lemma.strn like ?")
		resv = append(resv, q.LemmaLike)
	}
	if trm(q.LemmaRegexp) != "" {
		reses = append(reses, "Lemma.strn REGEXP ?")
		resv = append(resv, q.LemmaRegexp)
	}

	if trm(q.ReadingLike) != "" {
		reses = append(reses, "Lemma.reading like ?")
		resv = append(resv, q.ReadingLike)
	}
	if trm(q.ReadingRegexp) != "" {
		reses = append(reses, "Lemma.reading REGEXP ?")
		resv = append(resv, q.ReadingRegexp)
	}
	if trm(q.ParadigmLike) != "" {
		reses = append(reses, "Lemma.paradigm like ?")
		resv = append(resv, q.ParadigmLike)
	}
	if trm(q.ParadigmRegexp) != "" {
		reses = append(reses, "Lemma.paradigm REGEXP ?")
		resv = append(resv, q.ParadigmRegexp)
	}

	res := strings.Join(reses, " AND ")

	res += " and Lemma.id = Lemma2Entry.lemmaId and Entry.id = Lemma2Entry.entryId "

	return res, resv
}

func transcriptions(q Query) (string, []interface{}) {

	var reses []string
	var resv []interface{}

	if trm(q.TranscriptionLike) != "" {
		reses = append(reses, "Transcription.strn LIKE ?")
		resv = append(resv, q.TranscriptionLike)
	}

	if trm(q.TranscriptionRegexp) != "" {
		reses = append(reses, "Transcription.strn REGEXP ?")
		resv = append(resv, q.TranscriptionRegexp)
	}

	res := strings.Join(reses, " AND ")
	return res, resv
}

func entryStatuses(q Query) (string, []interface{}) {
	var res string
	var resv []interface{}

	if len(q.EntryStatus) == 0 {
		return res, resv
	}

	// Lexicon selection should already have been taken care of
	//res += "Entry.lexiconid = Lexicon.id AND Lexicon.id in " + nQs(len(q.Lexicons))
	res += "Entry.id = EntryStatus.entryId AND EntryStatus.current = 1 AND EntryStatus.name in " + nQs(len(q.EntryStatus))
	for _, es := range q.EntryStatus {
		resv = append(resv, es)
	}

	return res, resv
}

func users(q Query) (string, []interface{}) {
	var res string
	var resv []interface{}

	if len(q.Users) == 0 {
		return res, resv
	}

	// Lexicon selection should already have been taken care of
	//res += "Entry.lexiconid = Lexicon.id AND Lexicon.id in " + nQs(len(q.Lexicons))
	res += "Entry.id = EntryStatus.entryId AND EntryStatus.current = 1 AND EntryStatus.source in " + nQs(len(q.Users))
	for _, es := range q.Users {
		resv = append(resv, es)
	}

	return res, resv
}

func entryTag(q Query) (string, []interface{}) {

	var res []string
	var resv []interface{}

	if q.MultipleTags {
		res = append(res, " Entry.strn in (select distinct Entry.strn from EntryTag, Entry where EntryTag.entryId = Entry.id group by EntryTag.wordForm having count(EntryTag.wordForm)> 1)")
		//resv = append(resv, q.MultipleTags)
	}

	if q.TagLike != "" {
		res = append(res, " Entry.id = EntryTag.entryId AND EntryTag.tag like ? ")
		resv = append(resv, q.TagLike)
	}

	return strings.Join(res, " AND "), resv
}

func comments(q Query) (string, []interface{}) {
	var res []string
	var resv []interface{}

	if q.CommentLabelLike != "" {
		res = append(res, " Entry.id = EntryComment.entryId AND EntryComment.label like ? ")
		resv = append(resv, q.CommentLabelLike)
	}

	if q.CommentSourceLike != "" {
		res = append(res, " Entry.id = EntryComment.entryId AND EntryComment.source like ? ")
		resv = append(resv, q.CommentSourceLike)
	}

	if q.CommentLike != "" {
		res = append(res, " Entry.id = EntryComment.entryId AND EntryComment.comment like ? ")
		resv = append(resv, q.CommentLike)
	}

	return strings.Join(res, " AND "), resv
}

func validations(q Query) (string, []interface{}) {
	var res []string
	var resv []interface{}

	if q.ValidationRuleLike != "" {
		res = append(res, " Entry.id = EntryValidation.entryId AND lower(EntryValidation.name) like lower(?) ")
		resv = append(resv, q.ValidationRuleLike)
	}
	if q.ValidationLevelLike != "" {
		res = append(res, " Entry.id = EntryValidation.entryId AND lower(EntryValidation.level) like lower(?) ")
		resv = append(resv, q.ValidationLevelLike)
	}

	return strings.Join(res, " AND "), resv
}

func filter(ss []string, f func(string) bool) []string {
	var res []string
	for i, s := range ss {
		if f(s) {
			res = append(res, ss[i])
		}
	}
	return res
}

// RemoveEmptyStrings does that
func RemoveEmptyStrings(ss []string) []string {
	return filter(ss, func(s string) bool { return strings.TrimSpace(s) != "" })
}

// ToLower lower-cases its input strings
func ToLower(ss []string) []string {
	res := make([]string, len(ss))
	for i, v := range ss {
		res[i] = strings.ToLower(v)
	}
	return res
}

// This is not sane.

const baseSQLFrom = `FROM (Lexicon, Entry, Transcription)
LEFT JOIN Lemma2Entry ON Lemma2Entry.entryId = Entry.id 
LEFT JOIN Lemma ON Lemma.id = Lemma2Entry.lemmaid
LEFT JOIN EntryTag ON EntryTag.entryId = Entry.id
LEFT JOIN EntryStatus ON EntryStatus.entryId = Entry.id AND EntryStatus.current = 1
LEFT JOIN EntryValidation ON EntryValidation.entryId = Entry.id 
LEFT JOIN EntryComment ON EntryComment.entryId = Entry.id 
WHERE Entry.id = Transcription.entryId AND Entry.lexiconId = Lexicon.id` // Entry.lexiconid = Lexicon.id needed when no single input lexicon ID is given
// AND Lexicon.id = ? ORDER BY Entry.id, Transcription.id ASC`

// Queries db for all entries with transcriptions and optional lemma forms.
var baseSQLSelect = "SELECT Lexicon.name, Entry.id, Entry.strn, Entry.language, Entry.partOfSpeech, Entry.morphology, Entry.wordParts, Entry.preferred, Transcription.id, Transcription.entryId, Transcription.strn, Transcription.language, Transcription.sources, Lemma.id, Lemma.strn, Lemma.reading, Lemma.paradigm, EntryTag.tag, EntryStatus.id, EntryStatus.name, EntryStatus.source, EntryStatus.timestamp, EntryStatus.current, EntryValidation.id, EntryValidation.level, EntryValidation.name, EntryValidation.message, EntryValidation.timestamp, EntryComment.id, EntryComment.label, EntryComment.source, EntryComment.comment " + baseSQLFrom

//var baseSQLCount = `SELECT count(distinct Entry.id) ` + baseSQLFrom

var baseSQLSelectIds = `SELECT distinct Entry.id ` + baseSQLFrom

// sqlStmt container class for prepared sql statement:
// sql is a plain sql string with selects and '?' for arguments to be populated
// values is a range of values corresponding to the '?' arguments in the sql string
type sqlStmt struct {
	sql    string
	values []interface{}
}

func appendQuery(sql string, lexNames []lex.LexName, q Query) (string, []interface{}) {
	var args []interface{}

	// Query.Lexicons
	l, lv := lexicons(lexNames)
	args = append(args, lv...)
	// Query.Words, Query.WordsLike, Query.PartOfSpeechLike, Query.WordsRegexp, Query.PartOfSpeechRegexp
	w, wv := words(lexNames, q)
	args = append(args, wv...)
	// Query.Lemmas, Query.LemmaLike, Query.ReadingLike, Query.ParadigmLike, Query.LemmaRegexp, Query.ReadingRegexp, Query.ParadigmRegexp
	le, lev := lemmas(q)
	args = append(args, lev...)
	// Query.TranscriptionLike, Query.TranscriptionRegexp
	t, tv := transcriptions(q) // V2 simply returns 'transkription.strn like ?' + param value
	args = append(args, tv...)

	// Query.EntryStatus
	es, esv := entryStatuses(q)
	args = append(args, esv...)

	// Query.Users
	us, usv := users(q)
	args = append(args, usv...)

	// Query.TagLike
	tl, tlv := entryTag(q)
	args = append(args, tlv...)

	cl, clv := comments(q)
	args = append(args, clv...)

	vl, vlv := validations(q)
	args = append(args, vlv...)

	// HasEntryValidation doesn't take any argument
	ev := ""
	if q.HasEntryValidation {
		ev = "EntryValidation.entryId = Entry.id"
	}

	// puts together pieces of sql created above with " and " in between
	qRes := strings.TrimSpace(strings.Join(RemoveEmptyStrings([]string{l, w, le, t, es, us, tl, cl, vl, ev}), " AND "))
	if qRes != "" {
		sql += " AND " + qRes
	}
	// log.Printf("DEBUG QUERY %#v", q)
	// log.Printf("DEBUG QUERY RESULT %s\n\n", sql)
	return sql, args
}

// SelectEntriesSQL creates a SQL query string based on the values of
// a Query struct instance, along with a slice of values,
// corresponding to the params to be set (the '?':s of the query)
func selectEntriesSQL(lexNames []lex.LexName, q Query) sqlStmt {
	sqlQuery, args := appendQuery(baseSQLSelect, lexNames, q)

	// sort by id to make sql rows -> Entry simpler
	sqlQuery += " ORDER BY Entry.id, Transcription.id"

	// When both PageLength and Page values are zero, no page limit is used
	// This is useful for example when exporting a complete lexicon
	if q.PageLength > 0 || q.Page > 0 {
		sqlQuery += " LIMIT " + strconv.FormatInt(q.PageLength, 10) + " OFFSET " + strconv.FormatInt(q.PageLength*q.Page, 10)
	}
	return sqlStmt{sql: sqlQuery, values: args}

}

// SelectEntryIdsSQL creates a SQL query string based on the values of
// a Query struct instance, along with a slice of values,
// corresponding to the params to be set (the '?':s of the query)
func selectEntryIdsSQL(lexNames []lex.LexName, q Query) sqlStmt {
	sqlQuery, args := appendQuery(baseSQLSelectIds, lexNames, q)
	return sqlStmt{sql: sqlQuery, values: args}
}

// CountEntriesSQL creates a SQL query string based on the values of
// a Query struct instance, along with a slice of values,
// corresponding to the params to be set (the '?':s of the query)
/*
func countEntriesSQL(lexNames []lex.LexName, q Query) sqlStmt {
	sqlQuery, args := appendQuery(baseSQLCount, lexNames, q)
	return sqlStmt{sql: sqlQuery, values: args}
}
*/

// // entriesFromIdsSelect builds an sql select and returns it along with slice of matching id values
// func entriesFromIdsSelect(ids []int64) (string, []interface{}) {
// 	res := ""
// 	resv := convI(ids)
// 	qs := nQs(len(ids))
// 	// TODO assumes that every Entry has at least one transcription
// 	res += "select Entry.id, Entry.lexiconid, Entry.strn, Entry.language, Entry.partofspeech, Entry.wordparts, "
// 	res += "transcription.id, transcription.entryid, transcription.strn, transcription.language "
// 	res += "from Lexicon, entry, transcription "
// 	res += "where Lexicon.id = Entry.lexiconid "
// 	res += "and Entry.id = transcription.entryId "
// 	res += "and Entry.id in " + qs
// 	// res += " order by Entry.strn asc" // TODO ???
// 	return res, resv
// }
