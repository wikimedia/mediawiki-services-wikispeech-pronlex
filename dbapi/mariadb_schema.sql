-- TODO: In the Sqlite schema, TEXT type columns were used frequently, just to not have to decide the length.
-- In Mysql/MariDB, you cannot index a TEXT column without a length. Currently, such indecies are set to (255).
-- Probably, we should change the indexed TEXT columns to VARCHAR or maybe FULLTEXT index.
-- This also creates problems for multi-column indecies, since they have a max capacity.

-- CASE: It appears that MariaDB is case sensitive when it comes to table namnes, etc.


-- TODO: Remove!
DROP TABLE EntryTag, Entry, Lemma, Lexicon; 


-- To keep track of the version of this schema
-- CREATE TABLE SchemaVersion (name text not null);
-- INSERT INTO SchemaVersion VALUES (` + SchemaVersion + `);
-- Each lexical entry belongs to a lexicon.
-- The Lexicon table defines a lexicon through a unique name, along with the name a of symbol set and a locale
CREATE TABLE Lexicon (
    name varchar(128) not null,
    symbolSetName varchar(128) not null,
    locale varchar(128) not null,
    id integer not null primary key auto_increment
  );
CREATE UNIQUE INDEX name ON Lexicon (name);
CREATE UNIQUE INDEX namesymset ON Lexicon (name, symbolSetName);
-- Symbol set handling moved to file based solution
-- A symbol set is the definition of allowed symbols in a lexicons phonetical transcriptions
-- CREATE TABLE Symbolset (
    -- description varchar(128),
    -- description text,
    -- symbol varchar(128) not null,
    -- id integer not null primary key auto_increment,
    -- category varchar(128) not null,
    -- lexiconId integer not null,
    -- ipa varchar(128)
--   );
-- CREATE INDEX idx37380686 on Symbolset (symbol);
-- CREATE UNIQUE INDEX idx8bc90a52 on Symbolset (lexiconId,symbol);
-- Lemma forms, or stems, are uninflected (theoretical, one might say) forms of words
CREATE TABLE Lemma (
    reading varchar(128) not null,
    id integer not null primary key auto_increment,
    paradigm varchar(128),
    -- strn varchar(128) not null
    strn text not null
  );
CREATE INDEX reading on Lemma (reading);
CREATE INDEX paradigm on Lemma (paradigm);
CREATE INDEX strn on Lemma (strn(255));
CREATE INDEX lemidstrn on Lemma (id, strn(255));
-- TODO: NB: strn length is set to 128 since 255 as used elswhere is too
-- long in this multi-column index.
CREATE UNIQUE INDEX strnreading on Lemma (strn(128),reading);
--CREATE TABLE SurfaceForm (
--    id integer not null primary key auto_increment,
--    strn varchar(128) not null
--  );
--CREATE UNIQUE INDEX idx35390652 on SurfaceForm (strn);
-- The actual lexical entries live in this table.
-- Each entry is linked to a single lexicon, and may have one or more 
-- phonetic transcriptions, found in their own table.
CREATE TABLE Entry (
    -- wordParts varchar(128),
    wordParts text,
    label varchar(128), -- TODO What's this?!
    id integer not null primary key auto_increment,
    language varchar(128) not null,
    -- strn varchar(128) not null,
    strn text not null,
    lexiconId integer not null,
    partOfSpeech varchar(128),
    morphology varchar(128),
    preferred integer not null default 0, -- TODO Why doesn't it work when changing integer -> boolean? 
foreign key (lexiconId) references Lexicon(id));
CREATE INDEX language on Entry (language);
CREATE INDEX strn on Entry (strn(255));
CREATE INDEX lexiconid ON Entry (lexiconId);
CREATE INDEX entrypref ON Entry (preferred);
CREATE INDEX strnlangue on Entry (strn(255),language);
CREATE INDEX estrnpref on Entry (strn(255),preferred);
CREATE INDEX idid on Entry (id, lexiconId);
-- CREATE TABLE Tag (
--     strn text not null,
--     id integer not null primary key auto_increment,
-- );
-- CREATE UNIQUE INDEX tagindex ON Tag (strn);
-- Entry tag is a string used to distinguish between homographs.
-- Unique for an entry of a specific word form, but not for different
-- word forms. NOTE: This can be further normalized into a separate Tag
-- table, for reusable tags.
CREATE TABLE EntryTag (
    -- id integer not null primary key auto_increment,
    entryId integer not null,
    tag text not null,
    wordForm text, -- not null,
    FOREIGN KEY (entryId) REFERENCES Entry(id) ON DELETE CASCADE
);
-- A single tag per entry
CREATE UNIQUE INDEX tageid ON EntryTag(entryId);

-- TODO: NB: tag and wordForm length is set to 128 since 255 as used elswhere is too
-- long in this multi-column index.
CREATE UNIQUE INDEX tagentwf ON EntryTag(tag(128), wordForm(128));
-- Pick the entry word form from the Entry table
CREATE TRIGGER entryTagTrigger AFTER INSERT ON entryTag
   BEGIN
     UPDATE EntryTag SET wordForm = (select strn from Entry where id = entryId) WHERE EntryTag.entryId = NEW.entryId;
   END;
CREATE TRIGGER entryTagTrigger2 AFTER UPDATE ON entryTag
   BEGIN
     UPDATE EntryTag SET wordForm = (select strn from Entry where id = entryId) WHERE EntryTag.entryId = NEW.entryId;
   END;
CREATE TABLE EntryComment (
    id integer not null primary key auto_increment,
    entryId integer not null,
    source text,
    label text not null,
    comment text, -- not null,
    -- Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
    FOREIGN KEY (entryId) REFERENCES Entry(id) ON DELETE CASCADE
);
CREATE INDEX cmtlabelndx ON EntryComment(label(255)); 
CREATE INDEX cmtsrcndx ON EntryComment(source(255)); 
-- Validiation results of entries
CREATE TABLE EntryValidation (
    id integer not null primary key auto_increment,
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
    id integer not null primary key auto_increment,
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
    id integer not null primary key auto_increment,
    language varchar(128) not null,
    -- strn varchar(128) not null,
    strn text not null,
    sources TEXT not null,
foreign key (entryId) references Entry(id) on delete cascade);
CREATE INDEX traeid ON Transcription (entryId);
CREATE INDEX idtraeid ON Transcription (id, entryId);
-- CREATE TABLE TranscriptionStatus (
--    name varchar(128) not null,
--    source varchar(128) not null,
--    timestamp timestamp not null,
--    transcriptionId integer not null,
--    id integer not null primary key auto_increment,
-- foreign key (transcriptionId) references Transcription(id) on delete cascade);
-- CREATE INDEX nizze ON TranscriptionStatus (transcriptionId); 
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
-- CREATE TABLE SurfaceForm2Entry (
--    entryId bigint not null,
--    surfaceFormId bigint not null,
-- unique(surfaceFormId,entryId));
-- Triggers to ensure only one preferred = 1 per orthographic word
-- When a new entry is added, where preferred is not 0, all other entries for 
-- the same orthographic word (entry.strn), will have the preferred field set to 0.
CREATE TRIGGER insertPref BEFORE INSERT ON Entry
  BEGIN
    UPDATE Entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
  END;
CREATE TRIGGER updatePref BEFORE UPDATE ON Entry
  BEGIN
    UPDATE Entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
  END;
-- Triggers to ensure that there are only one entry status per entry
CREATE TRIGGER insertEntryStatus BEFORE INSERT ON EntryStatus
  BEGIN 
    UPDATE EntryStatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  END;
 CREATE TRIGGER updateEntryStatus BEFORE UPDATE ON EntryStatus
  BEGIN
    UPDATE EntryStatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  END;
