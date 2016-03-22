package dbapi

import (
	"strconv"
	"strings"
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
// entry.id from'
func tables(q Query) string {
	var res []string
	if len(q.Lexicons) > 0 {
		res = append(res, "lexicon")
	}
	//if len(q.Words) > 0 || q.WordLike != "" || q.PartOfSpeechLike != "" {
	res = append(res, "entry")
	//}

	if q.TranscriptionLike != "" {
		res = append(res, "transcription")
	}
	if len(q.Lemmas) > 0 || q.LemmaLike != "" || q.ReadingLike != "" || q.ParadigmLike != "" {
		res = append(res, "lemma, lemma2entry")
	}

	return strings.Join(res, ", ")
}

// lexicons returns a piece of sql matching the Querys list of
// lexicons and a slice of the db ids of the lexicons listed in a Query
func lexicons(q Query) (string, []interface{}) {
	var res string
	var resv []interface{}
	if q.Lexicons == nil || len(q.Lexicons) == 0 {
		return res, resv
	}

	res += "lexicon.id in " + nQs(len(q.Lexicons))

	lIds := make([]interface{}, len(q.Lexicons))
	for i, l := range q.Lexicons {
		lIds[i] = l.ID
	}

	resv = append(resv, lIds...)
	return res, resv
}

func words(q Query) (string, []interface{}) {
	var res string
	var resv []interface{}

	// If none of the following values are set, there are no
	// references to entries in the query. This should not make
	// sense, since we are building a query to look for entries,
	// but such a reference can be add by searching for other
	// things depending on entry (such as transcription)

	// You must study the Query struct to understand this

	if len(q.Words) == 0 && trm(q.WordLike) == "" && trm(q.WordRegexp) == "" && trm(q.PartOfSpeechLike) == "" && trm(q.PartOfSpeechRegexp) == "" && len(q.EntryIDs) == 0 {
		return res, resv
	} //else {
	if len(q.Words) > 0 {
		res += "entry.strn in " + nQs(len(q.Words))
		resv = append(resv, convS(ToLower(q.Words))...)
	}
	if len(q.EntryIDs) > 0 {
		res += "entry.id in " + nQs(len(q.EntryIDs))
		resv = append(resv, convI(q.EntryIDs)...)
	}
	if trm(q.WordLike) != "" {
		res += "entry.strn like ? "
		resv = append(resv, q.WordLike)
	}
	if trm(q.WordRegexp) != "" {
		res += "entry.strn REGEXP ? "
		resv = append(resv, q.WordRegexp)
	}

	if trm(q.PartOfSpeechLike) != "" {
		res += "entry.partofspeech like ? "
		resv = append(resv, q.PartOfSpeechLike)
	}
	if trm(q.PartOfSpeechRegexp) != "" {
		res += "entry.partofspeech REGEXP ? "
		resv = append(resv, q.PartOfSpeechRegexp)
	}

	//}

	if len(q.Lexicons) != 0 {
		res += " and entry.lexiconid = lexicon.id"
	}

	return res, resv
}

func lemmas(q Query) (string, []interface{}) {
	var res string
	var resv []interface{}

	if len(q.Lemmas) == 0 && trm(q.LemmaLike) == "" && trm(q.LemmaRegexp) == "" &&
		trm(q.ReadingLike) == "" && trm(q.ReadingRegexp) == "" &&
		trm(q.ParadigmLike) == "" && trm(q.ParadigmRegexp) == "" {
		return res, resv
	}
	if len(q.Lemmas) > 0 {
		res += "lemma.strn in " + nQs(len(q.Lemmas))
		resv = append(resv, convS(q.Lemmas)...)
	}
	if trm(q.LemmaLike) != "" {
		res += "lemma.strn like ? "
		resv = append(resv, q.LemmaLike)
	}
	if trm(q.LemmaRegexp) != "" {
		res += "lemma.strn REGEXP ? "
		resv = append(resv, q.LemmaRegexp)
	}

	if trm(q.ReadingLike) != "" {
		res += "lemma.reading like ? "
		resv = append(resv, q.ReadingLike)
	}
	if trm(q.ReadingRegexp) != "" {
		res += "lemma.reading REGEXP ? "
		resv = append(resv, q.ReadingRegexp)
	}
	if trm(q.ParadigmLike) != "" {
		res += "lemma.paradigm like ? "
		resv = append(resv, q.ParadigmLike)
	}
	if trm(q.ParadigmRegexp) != "" {
		res += "lemma.paradigm REGEXP ? "
		resv = append(resv, q.ParadigmRegexp)
	}

	res += " and lemma.id = lemma2entry.lemmaid and entry.id = lemma2entry.entryid "

	return res, resv
}

func transcriptions(q Query) (string, []interface{}) {

	var res string
	var resv []interface{}

	if trm(q.TranscriptionLike) != "" {
		res += "transcription.strn LIKE ? "
		resv = append(resv, q.ParadigmLike)
	}

	if trm(q.ParadigmRegexp) != "" {
		res += "transcription.strn REGEXP ? "
		resv = append(resv, q.ParadigmRegexp)
	}

	return res, resv
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

// Queries db for all entries with transcriptions and optional lemma forms.
var baseSQL = `SELECT lexicon.id, entry.id, entry.strn, entry.language, entry.partofspeech, entry.wordparts, transcription.id, transcription.entryid, transcription.strn, transcription.language, lemma.id, lemma.strn, lemma.reading, lemma.paradigm, entrystatus.id, entrystatus.name, entrystatus.source, entrystatus.timestamp, entrystatus.current  
FROM lexicon, entry, transcription 
LEFT JOIN lemma2entry ON lemma2entry.entryid = entry.id 
LEFT JOIN lemma ON lemma.id = lemma2entry.lemmaid 
LEFT JOIN entrystatus ON entrystatus.entryid = entry.id AND entrystatus.current = 1 
WHERE lexicon.id = entry.lexiconid AND entry.id = transcription.entryid ` // AND lexicon.id = ? ORDER BY entry.id, transcription.id ASC`

// SelectEntriesSQL creates a SQL query string based on the values of
// a Query struct instance, along with a slice of values,
// corresponding to the params to be set (the '?':s of the query)
func SelectEntriesSQL(q Query) (string, []interface{}) {
	var sqlQuery string
	var args []interface{}

	sqlQuery += baseSQL

	// Query.Lexicons
	l, lv := lexicons(q)
	args = append(args, lv...)
	// Query.Words, Query.WordsLike, Query.PartOfSpeechLike, Query.WordsRegexp, Query.PartOfSpeechRegexp
	w, wv := words(q)
	args = append(args, wv...)
	// Query.Lemmas, Query.LemmaLike, Query.ReadingLike, Query.ParadigmLike, Query.LemmaRegexp, Query.ReadingRegexp, Query.ParadigmRegexp
	le, lev := lemmas(q)
	args = append(args, lev...)
	// Query.TranscriptionLike, Query.TranscriptionRegexp
	t, tv := transcriptions(q) // V2 simply returns 'transkription.strn like ?' + param value
	args = append(args, tv...)

	// puts together pieces of sql created above with " and " in between
	qRes := strings.TrimSpace(strings.Join(RemoveEmptyStrings([]string{l, w, le, t}), " AND "))
	if "" != qRes {
		sqlQuery += " AND " + qRes
	}

	// sort by id to make sql rows -> Entry simpler
	sqlQuery += " ORDER BY entry.id, transcription.id"

	// When both PageLenth and Page values are zero, no page limit is used
	// This is useful for example when exporting a complete lexicon
	if q.PageLength > 0 || q.Page > 0 {
		sqlQuery += " LIMIT " + strconv.FormatInt(q.PageLength, 10) + " OFFSET " + strconv.FormatInt(q.PageLength*q.Page, 10)
	}
	return sqlQuery, args

}

// entriesFromIdsSelect builds an sql select and returns it along with slice of matching id values
func entriesFromIdsSelect(ids []int64) (string, []interface{}) {
	res := ""
	resv := convI(ids)
	qs := nQs(len(ids))
	// TODO assumes that every Entry has at least one transcription
	res += "select entry.id, entry.lexiconid, entry.strn, entry.language, entry.partofspeech, entry.wordparts, "
	res += "transcription.id, transcription.entryid, transcription.strn, transcription.language "
	res += "from lexicon, entry, transcription "
	res += "where lexicon.id = entry.lexiconid "
	res += "and entry.id = transcription.entryId "
	res += "and entry.id in " + qs
	// res += " order by entry.strn asc" // TODO ???
	return res, resv
}
