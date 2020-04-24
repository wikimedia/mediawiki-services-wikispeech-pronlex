   
   
     $ go run importLex.go 
          
     USAGE:
      importLex [FLAGS]
     
     FLAGS:
       -createdb
         	create db if it doesn't exist
       -db_engine string
         	db engine (sqlite or mariadb) (default "sqlite")
       -db_location string
         	db location (folder for sqlite; address for mariadb)
       -db_name string
         	db name
       -force
         	force loading of lexicon even if the symbolset is undefined (default: false)
       -help
         	print help message
       -lex_file string
         	lexicon file
       -lex_name string
         	lexicon name
       -locale string
         	lexicon locale
       -quiet
         	mute information logging (default: false)
       -symbolset string
         	lexicon symbolset file
       -validate
         	validate each entry, and save the validation in the database (default: false)
     
     SAMPLE INVOCATION:
       importLex -db_engine mariadb -db_location 'speechoid:@tcp(127.0.0.1:3306)' -lex_name sv-se.nst -locale sv_SE -lex_file [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8.gz -db_name svtest -symbolset [SYMBOLSET FOLDER]/sv-se_ws-sampa.sym 
       importLex -db_engine sqlite -db_location ~/wikispeech/sqlite -lex_name sv-se.nst -locale sv_SE -lex_file [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8.gz -db_name svtest -symbolset [SYMBOLSET FOLDER]/sv-se_ws-sampa.sym 
     
