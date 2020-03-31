-- TODO: In the Sqlite schema, TEXT type columns were used frequently, just to not have to decide the length.
-- In Mysql/MariDB, you cannot index a TEXT column without a length. Currently, such indecies are set to (255).
-- Probably, we should change the indexed TEXT columns to VARCHAR or maybe FULLTEXT index.
-- This also creates problems for multi-column indecies, since they have a max capacity.

-- CASE: It appears that MariaDB is case sensitive when it comes to table namnes, etc.


-- TODO: Remove!
DROP TABLE IF EXISTS SchemaVersion, Entry, Lexicon, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;
DROP TABLE IF EXISTS SchemaVersion, Entry, Lexicon, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;
DROP TABLE IF EXISTS SchemaVersion, Entry, Lexicon, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;
DROP TABLE IF EXISTS SchemaVersion, Entry, Lexicon, EntryComment, Lemma2Entry, Lemma, Transcription, EntryTag, EntryValidation, EntryStatus;

CREATE TABLE SchemaVersion (name text not null);
INSERT INTO SchemaVersion VALUES ('test_version_X');



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

CREATE TABLE Lemma (
    id integer not null primary key auto_increment,
    reading varchar(128) not null,
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

-- The actual lexical entries live in this table.
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
    foreign key fk_3  (lexiconId) references Lexicon(id));

CREATE INDEX language on Entry (language);
CREATE INDEX strn on Entry (strn(255));
CREATE INDEX lexiconid ON Entry (lexiconId);
CREATE INDEX entrypref ON Entry (preferred);
CREATE INDEX strnlangue on Entry (strn(255),language);
CREATE INDEX estrnpref on Entry (strn(255),preferred);
CREATE INDEX idid on Entry (id, lexiconId);

-- Entry tag is a string used to distinguish between homographs.
-- Unique for an entry of a specific word form, but not for different
-- word forms. NOTE: This can be further normalized into a separate Tag
-- table, for reusable tags.
CREATE TABLE EntryTag (
    -- id integer not null primary key auto_increment,
    entryId integer not null,
    tag text not null,
    wordForm text, -- not null,
    FOREIGN KEY fk_4 (entryId) REFERENCES Entry(id) ON DELETE CASCADE
);
-- A single tag per entry
CREATE UNIQUE INDEX tageid ON EntryTag(entryId);

-- TODO: NB: tag and wordForm length is set to 128 since 255 as used elswhere is too
-- long in this multi-column index.
CREATE UNIQUE INDEX tagentwf ON EntryTag(tag(128), wordForm(128));
-- Pick the entry word form from the Entry table
CREATE TRIGGER entryTagTrigger AFTER INSERT ON EntryTag
   FOR EACH ROW
     UPDATE EntryTag SET wordForm = (select strn from Entry where id = entryId) WHERE EntryTag.entryId = NEW.entryId;
   

CREATE TRIGGER entryTagTrigger2 AFTER UPDATE ON EntryTag
   FOR EACH ROW
     UPDATE EntryTag SET wordForm = (select strn from Entry where id = entryid) WHERE EntryTag.entryId = NEW.entryId;


CREATE TABLE EntryComment (
    id integer not null primary key auto_increment,
    entryId integer not null,
    source text,
    label text not null,
    comment text, -- not null,
    -- Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
    FOREIGN KEY fk_5 (entryId) REFERENCES Entry(id) ON DELETE CASCADE
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
    foreign key fk_6 (entryId) references Entry(id) on delete cascade);
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
    foreign key fk_7 (entryId) references Entry(id) on delete cascade);
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
    foreign key fk_8 (entryId) references Entry(id) on delete cascade);

CREATE INDEX traeid ON Transcription (entryId);
CREATE INDEX idtraeid ON Transcription (id, entryId);

-- As root, run SHOW ENGINE INNODB STATUS to get better error messages:
-- LATEST FOREIGN KEY ERROR
-- ------------------------
-- 2020-03-30 23:59:24 7f85a7621700 Error in foreign key constraint of table `test1`.`Lemma2Entry`:
--  FOREIGN KEY (`entryId`) REFERENCES Entry(`id`) ON DELETE CASCADE,
--     CONSTRAINT fk_2 FOREIGN KEY (`lemmaId`) REFERENCES Lemma(`id`) ON DELETE CASCADE):
-- Cannot find an index in the referenced table where the
-- referenced columns appear as the first columns, or column types
-- in the table and the referenced table do not match for constraint.
-- Note that the internal storage type of ENUM and SET changed in
-- tables created with >= InnoDB-4.1.12, and such columns in old tables
-- cannot be referenced by such columns in new tables.
-- See http://dev.mysql.com/doc/refman/5.6/en/innodb-foreign-key-constraints.html
-- for correct foreign key definition.
-- Create  table `test1`.`Lemma2Entry` with foreign key constraint failed. Field type or character set for column 'entryId' does not mach referenced column 'id' near ' FOREIGN KEY (`entryId`) REFERENCES Entry(`id`) ON DELETE CASCADE,
--     CONSTRAINT fk_2 FOREIGN KEY (`lemmaId`) REFERENCES Lemma(`id`) ON DELETE CASCADE)'.
-- ------------



-- Linking table between a lemma form and its different surface forms 
CREATE TABLE Lemma2Entry (
    entryId integer not null,
    lemmaId integer not null,
    unique(`lemmaId`,`entryId`),
    -- unique(`entryId`, `lemmaId`),
    FOREIGN KEY fk_1 (`entryId`) REFERENCES Entry(`id`) ON DELETE CASCADE,
    FOREIGN KEY fk_2 (`lemmaId`) REFERENCES Lemma(`id`) ON DELETE CASCADE);

--CREATE INDEX l2eind1 on Lemma2Entry (entryId);
CREATE INDEX l2eind2 on Lemma2Entry (lemmaId);
CREATE UNIQUE INDEX l2euind on Lemma2Entry (lemmaId,entryId);
CREATE UNIQUE INDEX idx46cf073d on Lemma2Entry (entryId);

-- Triggers to ensure only one preferred = 1 per orthographic word
-- When a new entry is added, where preferred is not 0, all other entries for 
-- the same orthographic word (entry.strn), will have the preferred field set to 0.
CREATE TRIGGER insertPref BEFORE INSERT ON Entry
  FOR EACH ROW
    UPDATE Entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
    
CREATE TRIGGER updatePref BEFORE UPDATE ON Entry
  FOR EACH ROW
    UPDATE Entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;

-- Triggers to ensure that there are only one entry status per entry
CREATE TRIGGER insertEntryStatus BEFORE INSERT ON EntryStatus
  FOR EACH ROW 
    UPDATE EntryStatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  
 CREATE TRIGGER updateEntryStatus BEFORE UPDATE ON EntryStatus
  FOR EACH ROW
    UPDATE EntryStatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  
