package dbapi

import "database/sql"

func LoadLexiconFile(db *sql.DB, logger Logger, lexiconName string, symbolSetName, lexiconFileName string) error {
	return nil

	// logger.Write(fmt.Sprintf("lexiconName: %v\n", lexiconName))
	// logger.Write(fmt.Sprintf("symbolSetName: %v\n", symbolSetName))
	// logger.Write(fmt.Sprintf("lexiconFileName: %v\n", lexiconFileName))

	// fh, err := os.Open(lexiconFileName)
	// if err != nil {
	// 	var msg = fmt.Sprintf("loadLexiconFileIntoDB failed to open file : %v", err)
	// 	logger.Write(msg)
	// 	return fmt.Errorf("loadLexiconFileIntoDB failed to open file : %v", err)
	// }

	// s := bufio.NewScanner(fh)

	// wsFmt, err := line.NewWS()
	// if err != nil {
	// 	var msg = fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
	// 	logger.Write(msg)
	// 	return fmt.Errorf("lexserver failed to instantiate lexicon line parser : %v", err)
	// }

	// TODO: Här är vi! Hur får vi fram lexiconID? /HL 20160826
	// lexicon := Lexicon{ID: lexiconID, Name: lexiconName, SymbolSetName: symbolSetName}

	// msg := fmt.Sprintf("Trying to load file: %s", lexiconFileName)
	// logger.Write(msg)

	// n := 0
	// var eBuf []lex.Entry
	// for s.Scan() {
	// 	if err := s.Err(); err != nil {
	// 		logger.Write(err.Error())
	// 		return fmt.Errorf("error when reading lines from lexicon file : %v", err)
	// 	}
	// 	l := s.Text()
	// 	e, err := wsFmt.ParseToEntry(l)
	// 	if err != nil {
	// 		logger.Write(err.Error())
	// 		return fmt.Errorf("error when parsing entry : %v", err)
	// 	}
	// 	eBuf = append(eBuf, e)
	// 	n++
	// 	if n%10000 == 0 {
	// 		_, err = InsertEntries(db, lexicon, eBuf)
	// 		if err != nil {
	// 			logger.Write(err.Error())
	// 			return fmt.Errorf("lexserver failed to insert entries : %v", err)
	// 		}
	// 		eBuf = make([]lex.Entry, 0)
	// 		msg2 := fmt.Sprintf("Lines so far: %d", n)
	// 		logger.Write(msg2)
	// 	}
	// }
	// InsertEntries(db, lexicon, eBuf) // flushing the buffer

	// _, err = db.Exec("ANALYZE")
	// if err != nil {
	// 	logger.Write(err.Error())
	// 	return fmt.Errorf("failed to exec analyze cmd to db : %v", err)
	// }

	// msg3 := fmt.Sprintf("Lines read:\t%d", n)
	// logger.Write(msg3)

	// if err := s.Err(); err != nil {
	// 	msg4 := fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
	// 	logger.Write(msg4)
	// 	return fmt.Errorf(msg4)
	// }

	// // TODO
	// return nil
}
