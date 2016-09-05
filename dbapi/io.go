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
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/vrules"
)

// Olika scenarion:
// Skapa
// Append
// Uppdatera

// func ImportSymbolSet(db *sql.DB, logger Logger, lexiconName, symbolSetName, lexiconFileName string) error {

// }

// ImportLexiconFile is intended for 'clean' imports. It doesn't check whether the words already exist and so on.
func ImportLexiconFile(db *sql.DB, logger Logger, lexiconName, lexiconFileName string, symbolSet symbolset.SymbolSet) []error {
	var errs []error

	ssRule := vrules.SymbolSetRule{symbolSet}
	ssName := symbolSet.Name

	logger.Write(fmt.Sprintf("lexiconName: %v", lexiconName))
	logger.Write(fmt.Sprintf("lexiconFileName: %v", lexiconFileName))
	logger.Write(fmt.Sprintf("symbolSetName: %v", ssName))

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

		for _, r := range ssRule.Validate(e) {
			logger.Write(r.String())
			errs = append(errs, fmt.Errorf("%v", r.String()))
			//return fmt.Errorf("%v", r)
			//panic(r) // shouldn't happen
		}

		eBuf = append(eBuf, e)
		n++
		if n%10000 == 0 {
			_, err = InsertEntries(db, lexicon, eBuf)
			if err != nil {
				var msg = fmt.Sprintf("lexserver failed to insert entries : %v", err)
				logger.Write(msg)
				errs = append(errs, fmt.Errorf("%v", msg))
				return errs
			}
			eBuf = make([]lex.Entry, 0)
			msg2 := fmt.Sprintf("Lines so far: %d", n)
			logger.Write(msg2)
		}
	}
	InsertEntries(db, lexicon, eBuf) // flushing the buffer

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
