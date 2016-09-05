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

// func ImportSymbolSet(db *sql.DB, logger Logger, lexiconName, symbolSetName, lexiconFileName string) error {

// }

// ImportLexiconFile is intended for 'clean' imports. It doesn't check whether the words already exist and so on.
func ImportLexiconFile(db *sql.DB, logger Logger, lexiconName, symbolSetName, lexiconFileName string) error {
	logger.Write(fmt.Sprintf("lexiconName: %v\n", lexiconName))
	logger.Write(fmt.Sprintf("symbolSetName: %v\n", symbolSetName))
	logger.Write(fmt.Sprintf("lexiconFileName: %v\n", lexiconFileName))

	fh, err := os.Open(lexiconFileName)
	if err != nil {
		var msg = fmt.Sprintf("ImportLexionFile failed to open file : %v", err)
		logger.Write(msg)
		return fmt.Errorf("ImportLexionFile failed to open file : %v", err)
	}

	var s *bufio.Scanner
	if strings.HasSuffix(lexiconFileName, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			var msg = fmt.Sprintf("ImportLexionFile failed to open gz reader : %v", err)
			logger.Write(msg)
			return fmt.Errorf("%v", msg)
		}
		s = bufio.NewScanner(gz)
	} else {
		s = bufio.NewScanner(fh)
	}

	wsFmt, err := line.NewWS()
	if err != nil {
		var msg = fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
		logger.Write(msg)
		return fmt.Errorf("lexserver failed to instantiate lexicon line parser : %v", err)
	}

	lexicon, err := GetLexicon(db, lexiconName)
	if err != nil {
		err = fmt.Errorf("lexserver failed to get lexicon id for lexicon: %s : %v", lexiconName, err)
		logger.Write(err.Error())
		return err
	}

	msg := fmt.Sprintf("Trying to load file: %s", lexiconFileName)
	logger.Write(msg)

	n := 0
	var eBuf []lex.Entry
	for s.Scan() {
		if err := s.Err(); err != nil {
			logger.Write(err.Error())
			return fmt.Errorf("error when reading lines from lexicon file : %v", err)
		}
		l := s.Text()
		e, err := wsFmt.ParseToEntry(l)
		if err != nil {
			logger.Write(err.Error())
			return fmt.Errorf("error when parsing entry : %v", err)
		}
		eBuf = append(eBuf, e)
		n++
		if n%10000 == 0 {
			_, err = InsertEntries(db, lexicon, eBuf)
			if err != nil {
				logger.Write(err.Error())
				return fmt.Errorf("lexserver failed to insert entries : %v", err)
			}
			eBuf = make([]lex.Entry, 0)
			msg2 := fmt.Sprintf("Lines so far: %d", n)
			logger.Write(msg2)
		}
	}
	InsertEntries(db, lexicon, eBuf) // flushing the buffer

	_, err = db.Exec("ANALYZE")
	if err != nil {
		logger.Write(err.Error())
		return fmt.Errorf("failed to exec analyze cmd to db : %v", err)
	}

	msg3 := fmt.Sprintf("Lines read:\t%d", n)
	logger.Write(msg3)

	if err := s.Err(); err != nil {
		msg4 := fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
		logger.Write(msg4)
		return fmt.Errorf(msg4)
	}

	// TODO
	return nil
}
