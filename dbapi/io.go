package dbapi

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/validation"
)

// ImportLexiconFile is intended for 'clean' imports. It doesn't check whether the words already exist and so on. It does not do any validation whatsoever of the transcriptions before they are added. If the validator parameter is initialized, each entry will be validated before import, and the validation result will be added to the db.
func ImportLexiconFile(db *sql.DB, logger Logger, lexiconName, lexiconFileName string, validator *validation.Validator) error {

	logger.Write(fmt.Sprintf("lexiconName: %v", lexiconName))
	logger.Write(fmt.Sprintf("lexiconFileName: %v", lexiconFileName))

	if _, err := os.Stat(lexiconFileName); os.IsNotExist(err) {
		var msg = fmt.Sprintf("ImportLexiconFile failed to open file : %v", err)
		logger.Write(msg)
		return fmt.Errorf("%v", msg)
	}

	fh, err := os.Open(lexiconFileName)
	defer fh.Close()
	if err != nil {
		var msg = fmt.Sprintf("ImportLexiconFile failed to open file : %v", err)
		logger.Write(msg)
		return fmt.Errorf("%v", msg)
	}

	var s *bufio.Scanner
	if strings.HasSuffix(lexiconFileName, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			var msg = fmt.Sprintf("ImportLexiconFile failed to open gz reader : %v", err)
			logger.Write(msg)
			return fmt.Errorf("%v", msg)
		}
		s = bufio.NewScanner(gz)
	} else {
		s = bufio.NewScanner(fh)
	}

	wsFmt, err := line.NewWS()
	if err != nil {
		var msg = fmt.Sprintf("ImportLexiconFile failed to instantiate lexicon line parser : %v", err)
		logger.Write(msg)
		return fmt.Errorf("%v", msg)
	}

	lexicon, err := GetLexicon(db, lexiconName)
	if err != nil {
		var msg = fmt.Sprintf("ImportLexiconFile failed to get lexicon id for lexicon: %s : %v", lexiconName, err)
		logger.Write(msg)
		return fmt.Errorf("%v", msg)
	}

	msg := fmt.Sprintf("Trying to load file: %s", lexiconFileName)
	logger.Write(msg)

	n := 0
	var eBuf []lex.Entry
	for s.Scan() {
		if err := s.Err(); err != nil {
			var msg = fmt.Sprintf("error when reading lines from lexicon file : %v", err)
			logger.Write(msg)
			return fmt.Errorf("%v", msg)
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
			return fmt.Errorf("%v", msg)
		}

		if validator != nil && validator.IsDefined() {
			validator.ValidateEntry(&e)
		}

		eBuf = append(eBuf, e)
		n++
		if n%10000 == 0 {
			msg2 := fmt.Sprintf("Inserting entries (total lines read: %d)  ...", n)
			logger.Write(msg2)
			_, err = InsertEntries(db, lexicon, eBuf)
			if err != nil {
				var msg = fmt.Sprintf("ImportLexiconFile failed to insert entries : %v", err)
				logger.Write(msg)
				return fmt.Errorf("%v", msg)
			} //else {
			//msg2 := fmt.Sprintf("Inserted entries (total lines read: %d)", n)
			logger.Write(msg2)
			//}
			eBuf = make([]lex.Entry, 0)
		}
		if n%logger.LogInterval() == 0 {
			msg2 := fmt.Sprintf("Lines read: %d", n)
			logger.Write(msg2)
		}
	}
	msg2 := fmt.Sprintf("Inserting entries (total lines read: %d)  ...", n)
	logger.Write(msg2)
	_, err = InsertEntries(db, lexicon, eBuf) // flushing the buffer
	if err != nil {
		var msg = fmt.Sprintf("ImportLexiconFile failed to insert entries : %v", err)
		logger.Write(msg)
		return fmt.Errorf("%v", msg)
	} //else {
	//msg2 := fmt.Sprintf("Inserted entries (total lines read: %d)", n)
	logger.Write(msg2)
	//}

	logger.Write("Finalizing import ... ")

	_, err = db.Exec("ANALYZE")
	if err != nil {
		var msg = fmt.Sprintf("failed to exec analyze cmd to db : %v", err)
		logger.Write(msg)
		return fmt.Errorf("%v", msg)
	}

	msg3 := fmt.Sprintf("Lines imported:\t%d", n)
	logger.Write(msg3)

	if err := s.Err(); err != nil {
		msg4 := fmt.Sprintf("ImportLexiconFile failed to instantiate lexicon line parser : %v", err)
		logger.Write(msg4)
		return fmt.Errorf("%v", msg4)
	}

	return nil
}

type PrintMode int

const (
	PrintAll PrintMode = iota
	PrintValid
	PrintInvalid
)

// ValidateLexiconFile validates the input file and prints any validation errors to the specified logger.
func ValidateLexiconFile(logger Logger, lexiconFileName string, validator *validation.Validator, printMode PrintMode) error {

	start := time.Now()

	var wg sync.WaitGroup
	log.Println(fmt.Sprintf("lexiconFileName: %v", lexiconFileName))

	if _, err := os.Stat(lexiconFileName); os.IsNotExist(err) {
		var msg = fmt.Sprintf("ValidateLexiconFile failed to open file : %v", err)
		log.Println(msg)
		return fmt.Errorf("%v", msg)
	}

	fh, err := os.Open(lexiconFileName)
	defer fh.Close()
	if err != nil {
		var msg = fmt.Sprintf("ValidateLexiconFile failed to open file : %v", err)
		log.Println(msg)
		return fmt.Errorf("%v", msg)
	}

	var s *bufio.Scanner
	if strings.HasSuffix(lexiconFileName, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			var msg = fmt.Sprintf("ValidateLexiconFile failed to open gz reader : %v", err)
			log.Println(msg)
			return fmt.Errorf("%v", msg)
		}
		s = bufio.NewScanner(gz)
	} else {
		s = bufio.NewScanner(fh)
	}

	wsFmt, err := line.NewWS()
	if err != nil {
		var msg = fmt.Sprintf("ValidateLexiconFile failed to instantiate lexicon line parser : %v", err)
		log.Println(msg)
		return fmt.Errorf("%v", msg)
	}

	msg := fmt.Sprintf("Trying to load file: %s", lexiconFileName)
	log.Println(msg)

	n := 0
	nPrinted := 0
	nValid := 0
	for s.Scan() {
		if err := s.Err(); err != nil {
			var msg = fmt.Sprintf("error when reading lines from lexicon file : %v", err)
			log.Println(msg)
			return fmt.Errorf("%v", msg)
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
			log.Println(msg)
			return fmt.Errorf("%v", msg)
		}

		wg.Add(1)
		go func(ee lex.Entry) {
			defer wg.Done()
			validator.ValidateEntry(&ee)
			isValid := (len(ee.EntryValidations) == 0)
			if isValid {
				nValid = nValid + 1
			}
			if printMode == PrintValid && isValid {
				logger.Write(l)
				nPrinted = nPrinted + 1
			} else if printMode == PrintAll {
				out := []string{l}
				nPrinted = nPrinted + 1
				for _, v := range ee.EntryValidations {
					out = append(out, fmt.Sprintf("#INVALID\t%#v", v.String()))
				}
				logger.Write(strings.Join(out, "\n"))
			} else if printMode == PrintInvalid && !isValid {
				out := []string{l}
				nPrinted = nPrinted + 1
				for _, v := range ee.EntryValidations {
					out = append(out, fmt.Sprintf("#INVALID\t%#v", v.String()))
				}
			}

			n++
			if n%logger.LogInterval() == 0 {
				msg2 := fmt.Sprintf("Lines read: %d", n)
				log.Println(msg2)
			}
		}(e)
	}
	wg.Wait()
	log.Printf("Lines read:\t%d", n)
	log.Printf("Lines printed:\t%d", nPrinted)
	log.Printf("Lines valid:\t%d", nValid)

	end := time.Now()
	log.Printf("dbapi/io.go ValidateLexiconFile took %v\n", end.Sub(start))

	if err := s.Err(); err != nil {
		msg = fmt.Sprintf("ValidateLexiconFile failed to instantiate lexicon line parser : %v", err)
		log.Println(msg)
		return fmt.Errorf("%v", msg)
	}

	return nil
}
