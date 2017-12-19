   
   
     $ go run importLex.go 
     
     USAGE:
      importLex <FLAGS> <DB FILE> <LEXICON NAME> <LOCALE> <LEXICON FILE> <SYMBOLSET NAME> <SYMBOLSET FOLDER>
     
     FLAGS:
        -validate bool  validate each entry, and save the validation in the database (default: false)
        -force    bool  force loading of lexicon even if the symbolset is undefined (default: false)
        -replace  bool  if the lexicon already exists, delete it before importing the new input data (default: false)
        -quiet    bool  mute information logging (default: false)
        -help     bool  print help message
     
     SAMPLE INVOCATION:
       importLex -validate pronlex.db sv-se.nst sv_SE [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8 sv-se_ws-sampa [SYMBOLSET FOLDER]
