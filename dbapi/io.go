package dbapi

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
)

// Olika scenarion:
// Skapa
// Append
// Uppdatera

// ImportLexiconFile is intended for 'clean' imports. It doesn't check whether the words already exist and so on. It does not do any validation whatsoever of the transcriptions.
func ImportLexiconFile(db *sql.DB, logger Logger, lexiconName, lexiconFileName string) []error {
	var errs []error

	logger.Write(fmt.Sprintf("lexiconName: %v", lexiconName))
	logger.Write(fmt.Sprintf("lexiconFileName: %v", lexiconFileName))

	fh, err := os.Open(lexiconFileName)
	defer fh.Close()
	if err != nil {
		var msg = fmt.Sprintf("ImportLexionFile failed to open file : %v", err)
		logger.Write(msg)
		errs = append(errs, fmt.Errorf("%v", msg))
		return errs
	}

	var s *bufio.Scanner
	if strings.HasSuffix(lexiconFileName, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			var msg = fmt.Sprintf("ImportLexionFile failed to open gz reader : %v", err)
			logger.Write(msg)
			errs = append(errs, fmt.Errorf("%v", msg))
			return errs
		}
		s = bufio.NewScanner(gz)
	} else {
		s = bufio.NewScanner(fh)
	}

	wsFmt, err := line.NewWS()
	if err != nil {
		var msg = fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
		logger.Write(msg)
		errs = append(errs, fmt.Errorf("%v", msg))
		return errs
	}

	lexicon, err := GetLexicon(db, lexiconName)
	if err != nil {
		var msg = fmt.Sprintf("lexserver failed to get lexicon id for lexicon: %s : %v", lexiconName, err)
		logger.Write(msg)
		errs = append(errs, fmt.Errorf("%v", msg))
		return errs
	}

	msg := fmt.Sprintf("Trying to load file: %s", lexiconFileName)
	logger.Write(msg)

	n := 0
	var eBuf []lex.Entry
	for s.Scan() {
		if err := s.Err(); err != nil {
			var msg = fmt.Sprintf("error when reading lines from lexicon file : %v", err)
			logger.Write(msg)
			errs = append(errs, fmt.Errorf("%v", msg))
			return errs
		}
		l := s.Text()

		if strings.HasPrefix(l, "#") {
			continue
		}
		if l == "" {
			continue
		}

		e, err := wsFmt.ParseToEntry(l)
		if err != nil {
			var msg = fmt.Sprintf("couldn't parse line to entry : %v", err)
			logger.Write(msg)
			errs = append(errs, fmt.Errorf("%v", msg))
			return errs
		}

		eBuf = append(eBuf, e)
		n++
		if n%10000 == 0 {
			msg2 := fmt.Sprintf("Inserting %d entries (total lines read: %d)  ...", len(eBuf), n)
			logger.Write(msg2)
			_, err = InsertEntries(db, lexicon, eBuf)
			if err != nil {
				var msg = fmt.Sprintf("lexserver failed to insert entries : %v", err)
				logger.Write(msg)
				errs = append(errs, fmt.Errorf("%v", msg))
				return errs
			} else {
				msg2 := fmt.Sprintf("Inserted %d entries (total lines read: %d)", len(eBuf), n)
				logger.Write(msg2)
			}
			eBuf = make([]lex.Entry, 0)
		}
		if n%logger.LogInterval() == 0 {
			msg2 := fmt.Sprintf("Lines read: %d", n)
			logger.Write(msg2)
		}
	}
	msg2 := fmt.Sprintf("Inserting %d entries (total lines read: %d)  ...", len(eBuf), n)
	logger.Write(msg2)
	_, err = InsertEntries(db, lexicon, eBuf) // flushing the buffer
	if err != nil {
		var msg = fmt.Sprintf("lexserver failed to insert entries : %v", err)
		logger.Write(msg)
		errs = append(errs, fmt.Errorf("%v", msg))
		return errs
	} else {
		msg2 := fmt.Sprintf("Inserted %d entries (total lines read: %d)", len(eBuf), n)
		logger.Write(msg2)
	}

	logger.Write("Finalizing ... ")

	_, err = db.Exec("ANALYZE")
	if err != nil {
		var msg = fmt.Sprintf("failed to exec analyze cmd to db : %v", err)
		logger.Write(msg)
		errs = append(errs, fmt.Errorf("%v", msg))
		return errs
	}

	msg3 := fmt.Sprintf("Lines read:\t%d", n)
	logger.Write(msg3)

	if err := s.Err(); err != nil {
		msg4 := fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
		logger.Write(msg4)
		errs = append(errs, fmt.Errorf("%v", msg4))
		return errs
	}

	return errs
}
