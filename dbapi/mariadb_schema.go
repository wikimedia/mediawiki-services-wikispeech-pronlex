package dbapi

// SchemaVersion defines the version of the schema structure. It is used for validating databases against the current version number. It will be updated manually when the structure of the schema/database is changed. Versions with the same prefix (e.g., 3 and 3.1) are compatible.
//const SchemaVersion = "3.1"

// TODO: SchemaVersion defined in schema.go

var MariaDBSchema = []string{

	`DROP TABLE IF EXISTS SchemaVersion, Lexion, Entry, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;`,

	`DROP TABLE IF EXISTS SchemaVersion, Entry, Lexicon, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;`,

	`DROP TABLE IF EXISTS SchemaVersion, Lexicon, Entry, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;`,

	`DROP TABLE IF EXISTS SchemaVersion, Entry, Lexicon, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;`,

	`CREATE TABLE SchemaVersion (name text not null);`,

	`INSERT INTO SchemaVersion VALUES (` + SchemaVersion + `);`,

	`CREATE TABLE Lexicon (
	    name varchar(128) not null,
	    symbolSetName varchar(128) not null,
	    locale varchar(128) not null,
	    id integer not null primary key auto_increment
	  );`,
	`CREATE UNIQUE INDEX name ON Lexicon (name);`,
	`CREATE UNIQUE INDEX namesymset ON Lexicon (name, symbolSetName);`,

	`CREATE TABLE Lemma (
	    id integer not null primary key auto_increment,
	    reading varchar(128) not null,
	    paradigm varchar(128),
	    -- strn varchar(128) not null
	    strn text not null
	  );`,
	`CREATE INDEX reading on Lemma (reading);`,
	`CREATE INDEX paradigm on Lemma (paradigm);`,
	`CREATE INDEX strn on Lemma (strn(255));`,
	`CREATE INDEX lemidstrn on Lemma (id, strn(255));`,
	`-- TODO: NB: strn length is set to 128 since 255 as used elswhere is too
	-- long in this multi-column index.
	CREATE UNIQUE INDEX strnreading on Lemma (strn(128),reading);`,

	`-- The actual lexical entries live in this table.
	-- Each entry is linked to a single lexicon, and may have one or more
	-- phonetic transcriptions, found in their own table.
	CREATE TABLE Entry (
	    -- wordParts varchar(128),
	    id integer not null primary key auto_increment,
	    wordParts text,
	    label varchar(128), -- TODO What's this?!
	    language varchar(128) not null,
	    -- strn varchar(128) not null,
	    strn text not null,
	    lexiconId integer not null,
	    partOfSpeech varchar(128),
	    morphology varchar(128),
	    preferred integer not null default 0, -- TODO Why doesn't it work when changing integer -> boolean?
	    foreign key fk_3  (lexiconId) references Lexicon(id));`,

	`CREATE INDEX language on Entry (language);`,
	`CREATE INDEX strn on Entry (strn(255));`,
	`CREATE INDEX lexiconId ON Entry (lexiconId);`,
	`CREATE INDEX entrypref ON Entry (preferred);`,
	`CREATE INDEX strnlangue on Entry (strn(255),language);`,
	`CREATE INDEX estrnpref on Entry (strn(255),preferred);`,
	`CREATE INDEX idid on Entry (id, lexiconId);`,

	`-- Entry tag is a string used to distinguish between homographs.
	-- Unique for an entry of a specific word form, but not for different
	-- word forms. NOTE: This can be further normalized into a separate Tag
	-- table, for reusable tags.
	CREATE TABLE EntryTag (
	    -- id integer not null primary key auto_increment,
	    entryId integer not null,
	    tag text not null,
	    wordForm text, -- not null,
	    FOREIGN KEY fk_4 (entryId) REFERENCES Entry(id) ON DELETE CASCADE
	);`,
	`-- A single tag per entry
	CREATE UNIQUE INDEX tageid ON EntryTag(entryId);`,

	`-- TODO: NB: tag and wordForm length is set to 128 since 255 as used elswhere is too
	-- long in this multi-column index.
	CREATE UNIQUE INDEX tagentwf ON EntryTag(tag(128), wordForm(128));`,
	`-- Pick the entry word form from the Entry table
	CREATE TRIGGER entryTagTrigger AFTER INSERT ON EntryTag
	   FOR EACH ROW
	     UPDATE EntryTag SET wordForm = (select strn from Entry where id = entryId) WHERE EntryTag.entryId = NEW.entryId;`,

	`CREATE TRIGGER entryTagTrigger2 AFTER UPDATE ON EntryTag
	   FOR EACH ROW
	     UPDATE EntryTag SET wordForm = (select strn from Entry where id = entryid) WHERE EntryTag.entryId = NEW.entryId;`,

	`CREATE TABLE EntryComment (
	    id integer not null primary key auto_increment,
	    entryId integer not null,
	    source text,
	    label text not null,
	    comment text, -- not null,
	    -- Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
	    FOREIGN KEY fk_5 (entryId) REFERENCES Entry(id) ON DELETE CASCADE
	);`,
	`CREATE INDEX cmtlabelndx ON EntryComment(label(255));`,
	`CREATE INDEX cmtsrcndx ON EntryComment(source(255));`,
	`-- Validiation results of entries
	CREATE TABLE EntryValidation (
	    id integer not null primary key auto_increment,
	    entryid integer not null,
	    level varchar(128) not null,
	    name varchar(128) not null,
	    -- message varchar(128) not null,
	    message text not null,
	    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
	    foreign key fk_6 (entryId) references Entry(id) on delete cascade);`,
	`CREATE INDEX evallev ON EntryValidation(level);`,
	`CREATE INDEX evalnam ON EntryValidation(name);`,
	`CREATE INDEX entvalEid ON EntryValidation(entryId);`,
	`CREATE INDEX identvalEid ON EntryValidation(id,entryId);`,
	`-- Status of entries
	CREATE TABLE EntryStatus (
	    name varchar(128) not null,
	    source varchar(128) not null,
	    entryId integer not null,
	    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
	    current boolean default 1 not null,
	    id integer not null primary key auto_increment,
	    UNIQUE(entryId,id),
	    foreign key fk_7 (entryId) references Entry(id) on delete cascade);`,
	`CREATE INDEX esn ON EntryStatus (name);`,
	`CREATE INDEX ess ON EntryStatus (source);`,
	`CREATE INDEX esc ON EntryStatus (current);`,
	`CREATE INDEX esceid ON EntryStatus (entryId);`,
	`CREATE INDEX entryidcurrent ON EntryStatus (entryId, current);`,
	`CREATE UNIQUE INDEX eseii ON EntryStatus  (id, entryId);`,
	`CREATE UNIQUE INDEX eseiicurr ON EntryStatus  (id, entryId, current);`,
	`CREATE UNIQUE INDEX idcurr ON EntryStatus  (id, current);`,
	`CREATE TABLE Transcription (
	    entryId integer not null,
	    preference int,
	    label varchar(128),
	    -- symbolSetCode varchar(128) not null,
	    id integer not null primary key auto_increment,
	    language varchar(128) not null,
	    -- strn varchar(128) not null,
	    strn text not null,
	    sources TEXT not null,
	    foreign key fk_8 (entryId) references Entry(id) on delete cascade);`,

	`CREATE INDEX traeid ON Transcription (entryId);`,
	`CREATE INDEX idtraeid ON Transcription (id, entryId);`,

	`-- Linking table between a lemma form and its different surface forms
	CREATE TABLE Lemma2Entry (
	    entryId integer not null,
	    lemmaId integer not null,
	    unique(lemmaId,entryId),
	    -- unique(entryId, lemmaId),
	    FOREIGN KEY fk_1 (entryId) REFERENCES Entry(id) ON DELETE CASCADE,
	    FOREIGN KEY fk_2 (lemmaId) REFERENCES Lemma(id) ON DELETE CASCADE);`,

	`CREATE INDEX l2eind2 on Lemma2Entry (lemmaId);`,
	`CREATE UNIQUE INDEX l2euind on Lemma2Entry (lemmaId,entryId);`,
	`CREATE UNIQUE INDEX idx46cf073d on Lemma2Entry (entryId);`,

	`-- Triggers to ensure only one preferred = 1 per orthographic word
	-- When a new entry is added, where preferred is not 0, all other entries for
	-- the same orthographic word (entry.strn), will have the preferred field set to 0.
	-- TODO: This doesn't work in Mysql/MariaDB?
        CREATE TRIGGER insertPref BEFORE INSERT ON Entry
	  FOR EACH ROW
          BEGIN
	    UPDATE Entry SET NEW.preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconId = NEW.lexiconId;
END;`,

	// `CREATE TRIGGER updatePref BEFORE UPDATE ON Entry
	//   FOR EACH ROW
	//     UPDATE Entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconId = NEW.lexiconId;`,

	`-- Triggers to ensure that there are only one entry status per entry
	CREATE TRIGGER insertEntryStatus BEFORE INSERT ON EntryStatus
	  FOR EACH ROW
	    UPDATE EntryStatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;`,

	`CREATE TRIGGER updateEntryStatus BEFORE UPDATE ON EntryStatus
	  FOR EACH ROW
	    UPDATE EntryStatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;`,
}
