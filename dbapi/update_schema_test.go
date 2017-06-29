package dbapi

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	//"github.com/stts-se/pronlex/dbapi"
)

var tstSchema1 = `
PRAGMA user_version = 0;
-- Each lexical entry belongs to a lexicon.
-- The Lexicon table defines a lexicon through a unique name, along with the name a of symbol set
CREATE TABLE Lexicon (
    name varchar(128) not null,
    symbolSetName varchar(128) not null,
    id integer not null primary key autoincrement
  );
CREATE UNIQUE INDEX idx1e0404a1 on Lexicon (name);
CREATE UNIQUE INDEX namesymset on Lexicon (name, symbolSetName);

-- Lemma forms, or stems, are uninflected (theoretical, one might say) forms of words
CREATE TABLE Lemma (
    reading varchar(128) not null,
    id integer not null primary key autoincrement,
    paradigm varchar(128),
    -- strn varchar(128) not null
    strn text not null
  );
CREATE INDEX idx21d604f4 on Lemma (reading);
CREATE INDEX idx273f055f on Lemma (paradigm);
CREATE INDEX idx149303e1 on Lemma (strn);
CREATE INDEX lemidstrn on Lemma (id, strn);
CREATE UNIQUE INDEX idx407206e8 on Lemma (strn,reading);
CREATE TABLE Entry (
    -- wordParts varchar(128),
    wordParts text,
    label varchar(128),
    id integer not null primary key autoincrement,
    language varchar(128) not null,
    -- strn varchar(128) not null,
    strn text not null,
    lexiconId integer not null,
    partOfSpeech varchar(128),
    morphology varchar(128),
    preferred integer not null default 0, -- TODO Why doesn't it work when changing integer -> boolean? 
foreign key (lexiconId) references Lexicon(id));
CREATE INDEX idx28d70584 on Entry (language);
CREATE INDEX idx15890407 on Entry (strn);
CREATE INDEX entrylexid ON Entry (lexiconId);
CREATE INDEX entrypref ON Entry (preferred);
CREATE INDEX idx4a250778 on Entry (strn,language);
CREATE INDEX estrnpref on Entry (strn,preferred);
CREATE INDEX idid on Entry (id, lexiconId);

-- Validiation results of entries
CREATE TABLE EntryValidation (
    id integer not null primary key autoincrement,
    entryid integer not null,
    level varchar(128) not null,
    name varchar(128) not null,
    -- message varchar(128) not null,
    message text not null,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
    foreign key (entryId) references Entry(id) on delete cascade);
CREATE INDEX evallev ON EntryValidation(level);
CREATE INDEX evalnam ON EntryValidation(name);
CREATE INDEX entvalEid ON EntryValidation(entryId); 
CREATE INDEX identvalEid ON EntryValidation(id,entryId); 

-- Status of entries
CREATE TABLE EntryStatus (
    name varchar(128) not null,
    source varchar(128) not null,
    entryId integer not null,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
    current boolean default 1 not null,
    id integer not null primary key autoincrement,
    UNIQUE(entryId,id),
    foreign key (entryId) references Entry(id) on delete cascade);
CREATE INDEX esn ON EntryStatus (name);
CREATE INDEX ess ON EntryStatus (source);
CREATE INDEX esc ON EntryStatus (current);
CREATE INDEX esceid ON EntryStatus (entryId);
CREATE INDEX entryidcurrent ON EntryStatus (entryId, current);
CREATE UNIQUE INDEX eseii ON EntryStatus  (id, entryId);
CREATE UNIQUE INDEX eseiicurr ON EntryStatus  (id, entryId, current);
CREATE UNIQUE INDEX idcurr ON EntryStatus  (id, current);

CREATE TABLE Transcription (
    entryId integer not null,
    preference int,
    label varchar(128),
    -- symbolSetCode varchar(128) not null,
    id integer not null primary key autoincrement,
    language varchar(128) not null,
    -- strn varchar(128) not null,
    strn text not null,
    sources TEXT not null,
foreign key (entryId) references Entry(id) on delete cascade);
CREATE INDEX traeid ON Transcription (entryId);
CREATE INDEX idtraeid ON Transcription (id, entryId);

-- Linking table between a lemma form and its different surface forms 
CREATE TABLE Lemma2Entry (
    entryId bigint not null,
    lemmaId bigint not null,
unique(lemmaId,entryId),
foreign key (entryId) references Entry(id) on delete cascade,
foreign key (lemmaId) references Lemma(id) on delete cascade);
--CREATE INDEX l2eind1 on Lemma2Entry (entryId);
CREATE INDEX l2eind2 on Lemma2Entry (lemmaId);
CREATE UNIQUE INDEX l2euind on Lemma2Entry (lemmaId,entryId);
CREATE UNIQUE INDEX idx46cf073d on Lemma2Entry (entryId);


-- Triggers to ensure only one preferred = 1 per orthographic word
-- When a new entry is added, where preferred is not 0, all other entries for 
-- the same orthographic word (entry.strn), will have the preferred field set to 0.
CREATE TRIGGER insertPref BEFORE INSERT ON ENTRY
  BEGIN
    UPDATE entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0; -- BROKEN! AND lexiconid = NEW.lexiconid;
  END;
CREATE TRIGGER updatePref BEFORE UPDATE ON ENTRY
  BEGIN
    UPDATE entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
  END;
`

var trigger1 = `CREATE TRIGGER insertPref BEFORE INSERT ON ENTRY
  BEGIN
    UPDATE entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
  END;`

var trigger2 = `-- Triggers to ensure that there are only one entry status per entry
CREATE TRIGGER insertEntryStatus BEFORE INSERT ON ENTRYSTATUS
  BEGIN 
    UPDATE entrystatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  END;
CREATE TRIGGER updateEntryStatus BEFORE UPDATE ON ENTRYSTATUS
  BEGIN
    UPDATE entrystatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  END;
`

func Test_UpdateSchema(t *testing.T) {

	// This could be put in some set-up function

	tstDBName := "/tmp/tztdb.db"

	if _, err := os.Stat(tstDBName); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "removing old test db '%s'\n", tstDBName)
		err := os.Remove(tstDBName)
		if err != nil {
			t.Errorf("Weltschmerz! : %v", err)
		}
	}

	db, err := sql.Open("sqlite3", tstDBName)
	defer db.Close()
	if err != nil {
		t.Errorf("This isn't happening! : %v", err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Errorf("Please...! : %v", err)
	}

	_, err = db.Exec(tstSchema1)
	if err != nil {
		t.Errorf("No one saw this one coming : %v", err)
	}

	if err := UpdateSchema(tstDBName); err != nil {
		t.Errorf("NOOOOOOOOOOOOOOOOOOOO: %v", err)
	}

	// rows, err := db.Query("SELECT * FROM entry limit 0") //db.Query("PRAGMA table_info('entry')")
	// defer rows.Close()

	// nms, err := listNamesOfTriggers(db) // defined in dbapi.go
	// if err != nil {
	// 	fmt.Printf("APNOS! %v\n", err)
	// }

	// for _, n := range nms {
	// 	fmt.Println(">>>> " + n)
	// }

	// for rows.Next() {
	// 	var id interface{}
	// 	var name interface{}
	// 	var c3 interface{}
	// 	var c4 interface{}
	// 	var c5 interface{}
	// 	var c6 interface{}

	// 	if err := rows.Scan(&id, &name, &c3, &c4, &c5...); err == nil {
	// 		fmt.Printf("%v %s\n", id, name)
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	// }
}
