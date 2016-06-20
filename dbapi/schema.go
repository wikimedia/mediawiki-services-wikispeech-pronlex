package dbapi

// Schema is a string containing the SQL definition of the lexicon database
var Schema = `
-- Each lexical entry belongs to a lexicon.
-- The Lexicon table defines a lexicon through a unique name, along with the name a of symbol set
CREATE TABLE Lexicon (
    name varchar(128) not null,
    symbolSetName varchar(128) not null,
    id integer not null primary key autoincrement
  );
CREATE UNIQUE INDEX idx1e0404a1 on Lexicon (name);

-- A symbol set is the definition of allowed symbols in a lexicons phonetical transcriptions
CREATE TABLE Symbolset (
    description varchar(128),
    symbol varchar(128) not null,
    id integer not null primary key autoincrement,
    category varchar(128) not null,
    lexiconId integer not null,
    ipa varchar(128)
  );
CREATE INDEX idx37380686 on Symbolset (symbol);
CREATE UNIQUE INDEX idx8bc90a52 on Symbolset (lexiconId,symbol);

-- Lemma forms, or stems, are uninflected (theoretical, one might say) forms of words
CREATE TABLE Lemma (
    reading varchar(128) not null,
    id integer not null primary key autoincrement,
    paradigm varchar(128),
    strn varchar(128) not null
  );
CREATE INDEX idx21d604f4 on Lemma (reading);
CREATE INDEX idx273f055f on Lemma (paradigm);
CREATE INDEX idx149303e1 on Lemma (strn);
CREATE UNIQUE INDEX idx407206e8 on Lemma (strn,reading);
--CREATE TABLE SurfaceForm (
--    id integer not null primary key autoincrement,
--    strn varchar(128) not null
--  );
--CREATE UNIQUE INDEX idx35390652 on SurfaceForm (strn);

-- The actual lexical entries live in this table.
-- Each entry is linked to a single lexicon, and may have one or more 
-- phonetic transcriptions, found in their own table.
CREATE TABLE Entry (
    wordParts varchar(128),
    label varchar(128),
    id integer not null primary key autoincrement,
    language varchar(128) not null,
    strn varchar(128) not null,
    lexiconId integer not null,
    partOfSpeech varchar(128),
foreign key (lexiconId) references Lexicon(id));
CREATE INDEX idx28d70584 on Entry (language);
CREATE INDEX idx15890407 on Entry (strn);
CREATE INDEX entrylexid ON Entry (lexiconId);
CREATE INDEX idx4a250778 on Entry (strn,language);


-- Validiation results of entries
CREATE TABLE EntryValidation (
    id integer not null primary key autoincrement,
    entryid integer not null,
    level varchar(128) not null,
    name varchar(128) not null,
    message varchar(128) not null,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
    foreign key (entryId) references Entry(id) on delete cascade);
CREATE INDEX evallev ON EntryValidation(level);
CREATE INDEX evalnam ON EntryValidation(name);
CREATE INDEX entvalEid ON EntryValidation(entryId); 

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
CREATE UNIQUE INDEX eseii ON EntryStatus  (entryId, id);

CREATE TABLE Transcription (
    entryId integer not null,
    preference int,
    label varchar(128),
    -- symbolSetCode varchar(128) not null,
    id integer not null primary key autoincrement,
    language varchar(128) not null,
    strn varchar(128) not null,
    sources TEXT not null,
foreign key (entryId) references Entry(id) on delete cascade);
CREATE INDEX traeid ON Transcription (entryId);

CREATE TABLE TranscriptionStatus (
    name varchar(128) not null,
    source varchar(128) not null,
    timestamp timestamp not null,
    transcriptionId integer not null,
    id integer not null primary key autoincrement,
foreign key (transcriptionId) references Transcription(id) on delete cascade);
CREATE INDEX nizze ON TranscriptionStatus (transcriptionId); 

-- Linking table between a lemma form and its different surface forms 
CREATE TABLE Lemma2Entry (
    entryId bigint not null,
    lemmaId bigint not null,
unique(lemmaId,entryId),
foreign key (entryId) references Entry(id) on delete cascade,
foreign key (lemmaId) references Lemma(id) on delete cascade);
CREATE INDEX l2eind1 on Lemma2Entry (entryId);
CREATE INDEX l2eind2 on Lemma2Entry (lemmaId);
CREATE UNIQUE INDEX l2euind on Lemma2Entry (lemmaId,entryId);
CREATE UNIQUE INDEX idx46cf073d on Lemma2Entry (entryId);

-- CREATE TABLE SurfaceForm2Entry (
--    entryId bigint not null,
--    surfaceFormId bigint not null,
-- unique(surfaceFormId,entryId));
`
