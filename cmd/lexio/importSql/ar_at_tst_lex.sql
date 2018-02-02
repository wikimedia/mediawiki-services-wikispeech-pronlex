PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE Lexicon (
    name varchar(128) not null,
    symbolSetName varchar(128) not null,
    locale varchar(128) not null,
    id integer not null primary key autoincrement
  );
INSERT INTO Lexicon VALUES('ar-test','ar_ws-sampa','ar_AR',1);
CREATE TABLE Lemma (
    reading varchar(128) not null,
    id integer not null primary key autoincrement,
    paradigm varchar(128),
    -- strn varchar(128) not null
    strn text not null
  );
CREATE TABLE Entry (
    -- wordParts varchar(128),
    wordParts text,
    label varchar(128), -- TODO What's this?!
    id integer not null primary key autoincrement,
    language varchar(128) not null,
    -- strn varchar(128) not null,
    strn text not null,
    lexiconId integer not null,
    partOfSpeech varchar(128),
    morphology varchar(128),
    preferred integer not null default 0, -- TODO Why doesn't it work when changing integer -> boolean? 
foreign key (lexiconId) references Lexicon(id));
INSERT INTO Entry VALUES('bob',NULL,1,'en','bob',1,'','',0);
INSERT INTO Entry VALUES('dylan',NULL,2,'en','dylan',1,'','',0);
INSERT INTO Entry VALUES('volvo',NULL,3,'en','volvo',1,'','',0);
INSERT INTO Entry VALUES('بوب',NULL,4,'ar','بوب',1,'','',0);
INSERT INTO Entry VALUES('ديلن',NULL,5,'ar','ديلن',1,'','',0);
INSERT INTO Entry VALUES('فولفو',NULL,6,'en','فولفو',1,'','',0);
CREATE TABLE EntryTag (
    -- id integer not null primary key autoincrement,
    entryId integer not null,
    tag text not null,
    wordForm text, -- not null,
    FOREIGN KEY (entryId) REFERENCES Entry(id) ON DELETE CASCADE
);
CREATE TABLE EntryValidation (
    id integer not null primary key autoincrement,
    entryid integer not null,
    level varchar(128) not null,
    name varchar(128) not null,
    -- message varchar(128) not null,
    message text not null,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
    foreign key (entryId) references Entry(id) on delete cascade);
CREATE TABLE EntryStatus (
    name varchar(128) not null,
    source varchar(128) not null,
    entryId integer not null,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP not null,
    current boolean default 1 not null,
    id integer not null primary key autoincrement,
    UNIQUE(entryId,id),
    foreign key (entryId) references Entry(id) on delete cascade);
INSERT INTO EntryStatus VALUES('imported','hb',1,'2017-11-28 10:12:15',1,1);
INSERT INTO EntryStatus VALUES('imported','hb',2,'2017-11-28 10:12:15',1,2);
INSERT INTO EntryStatus VALUES('imported','hb',3,'2017-11-28 10:12:15',1,3);
INSERT INTO EntryStatus VALUES('imported','hb',4,'2017-11-28 10:12:15',1,4);
INSERT INTO EntryStatus VALUES('imported','hb',5,'2017-11-28 10:12:15',1,5);
INSERT INTO EntryStatus VALUES('imported','hb',6,'2017-11-28 10:12:15',1,6);
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
INSERT INTO Transcription VALUES(1,NULL,NULL,1,'ar',''' b o b','hb');
INSERT INTO Transcription VALUES(2,NULL,NULL,2,'ar',''' d i . l a n','hb');
INSERT INTO Transcription VALUES(3,NULL,NULL,3,'ar',''' v o l . v o:','hb');
INSERT INTO Transcription VALUES(4,NULL,NULL,4,'ar',''' b o b','hb');
INSERT INTO Transcription VALUES(5,NULL,NULL,5,'ar',''' d i . l a n','hb');
INSERT INTO Transcription VALUES(6,NULL,NULL,6,'ar',':'' v o l . v o','hb');
CREATE TABLE Lemma2Entry (
    entryId bigint not null,
    lemmaId bigint not null,
unique(lemmaId,entryId),
foreign key (entryId) references Entry(id) on delete cascade,
foreign key (lemmaId) references Lemma(id) on delete cascade);
ANALYZE sqlite_master;
INSERT INTO sqlite_stat1 VALUES('Lexicon','namesymset','1 1 1');
INSERT INTO sqlite_stat1 VALUES('Lexicon','idx1e0404a1','1 1');
INSERT INTO sqlite_stat1 VALUES('Entry','idid','6 1 1');
INSERT INTO sqlite_stat1 VALUES('Entry','estrnpref','6 1 1');
INSERT INTO sqlite_stat1 VALUES('Entry','idx4a250778','6 1 1');
INSERT INTO sqlite_stat1 VALUES('Entry','entrypref','6 6');
INSERT INTO sqlite_stat1 VALUES('Entry','entrylexid','6 6');
INSERT INTO sqlite_stat1 VALUES('Entry','idx15890407','6 1');
INSERT INTO sqlite_stat1 VALUES('Entry','idx28d70584','6 3');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','idcurr','6 1 1');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','eseiicurr','6 1 1 1');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','eseii','6 1 1');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','entryidcurrent','6 1 1');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','esceid','6 1');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','esc','6 6');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','ess','6 6');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','esn','6 6');
INSERT INTO sqlite_stat1 VALUES('EntryStatus','sqlite_autoindex_EntryStatus_1','6 1 1');
INSERT INTO sqlite_stat1 VALUES('Transcription','idtraeid','6 1 1');
INSERT INTO sqlite_stat1 VALUES('Transcription','traeid','6 1');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('Lexicon',1);
INSERT INTO sqlite_sequence VALUES('Entry',6);
INSERT INTO sqlite_sequence VALUES('Transcription',6);
INSERT INTO sqlite_sequence VALUES('EntryStatus',6);
CREATE UNIQUE INDEX idx1e0404a1 on Lexicon (name);
CREATE UNIQUE INDEX namesymset on Lexicon (name, symbolSetName);
CREATE INDEX idx21d604f4 on Lemma (reading);
CREATE INDEX idx273f055f on Lemma (paradigm);
CREATE INDEX idx149303e1 on Lemma (strn);
CREATE INDEX lemidstrn on Lemma (id, strn);
CREATE UNIQUE INDEX idx407206e8 on Lemma (strn,reading);
CREATE INDEX idx28d70584 on Entry (language);
CREATE INDEX idx15890407 on Entry (strn);
CREATE INDEX entrylexid ON Entry (lexiconId);
CREATE INDEX entrypref ON Entry (preferred);
CREATE INDEX idx4a250778 on Entry (strn,language);
CREATE INDEX estrnpref on Entry (strn,preferred);
CREATE INDEX idid on Entry (id, lexiconId);
CREATE UNIQUE INDEX tageid ON EntryTag(entryId);
CREATE UNIQUE INDEX tagentwf ON EntryTag(tag, wordForm);
CREATE TRIGGER entryTagTrigger AFTER INSERT ON entryTag
   BEGIN
     UPDATE EntryTag SET wordForm = (select strn from entry where id = entryid) WHERE EntryTag.entryId = NEW.entryId;
   END;
CREATE TRIGGER entryTagTrigger2 AFTER UPDATE ON entryTag
   BEGIN
     UPDATE EntryTag SET wordForm = (select strn from entry where id = entryid) WHERE EntryTag.entryId = NEW.entryId;
   END;
CREATE INDEX evallev ON EntryValidation(level);
CREATE INDEX evalnam ON EntryValidation(name);
CREATE INDEX entvalEid ON EntryValidation(entryId);
CREATE INDEX identvalEid ON EntryValidation(id,entryId);
CREATE INDEX esn ON EntryStatus (name);
CREATE INDEX ess ON EntryStatus (source);
CREATE INDEX esc ON EntryStatus (current);
CREATE INDEX esceid ON EntryStatus (entryId);
CREATE INDEX entryidcurrent ON EntryStatus (entryId, current);
CREATE UNIQUE INDEX eseii ON EntryStatus  (id, entryId);
CREATE UNIQUE INDEX eseiicurr ON EntryStatus  (id, entryId, current);
CREATE UNIQUE INDEX idcurr ON EntryStatus  (id, current);
CREATE INDEX traeid ON Transcription (entryId);
CREATE INDEX idtraeid ON Transcription (id, entryId);
CREATE INDEX l2eind2 on Lemma2Entry (lemmaId);
CREATE UNIQUE INDEX l2euind on Lemma2Entry (lemmaId,entryId);
CREATE UNIQUE INDEX idx46cf073d on Lemma2Entry (entryId);
CREATE TRIGGER insertPref BEFORE INSERT ON ENTRY
  BEGIN
    UPDATE entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
  END;
CREATE TRIGGER updatePref BEFORE UPDATE ON ENTRY
  BEGIN
    UPDATE entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
  END;
CREATE TRIGGER insertEntryStatus BEFORE INSERT ON ENTRYSTATUS
  BEGIN 
    UPDATE entrystatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  END;
CREATE TRIGGER updateEntryStatus BEFORE UPDATE ON ENTRYSTATUS
  BEGIN
    UPDATE entrystatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
  END;
COMMIT;
