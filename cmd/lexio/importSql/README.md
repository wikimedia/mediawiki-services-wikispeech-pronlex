   
   
     $ go run importSql.go 
     
     USAGE:
      importSql <SQL DUMP FILE> <NEW DB FILE>

      <SQL DUMP FILE> - sql dump of a lexicon database (.sql or .sql.gz)
      <NEW DB FILE>   - new (non-existing) db file to import into (<DBNAME>.db)
     
     SAMPLE INVOCATION:
       importSql [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8.sql.gz sv_se_nst_lex.db
