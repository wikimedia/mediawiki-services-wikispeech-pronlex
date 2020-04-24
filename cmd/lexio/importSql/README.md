
     $ go run importSql.go  -h   
     
     USAGE:
           importSql [FLAGS] <SQL DUMP FILE>
     
           <SQL DUMP FILE> - sql dump of a lexicon database (.sql or .sql.gz)
          
          SAMPLE INVOCATION:
            importSql go run . -db_engine mariadb -mariadb_user speechoid -mariadb_host 127.0.0.1 -db_name testfest swe030224NST.pron-ws.utf8.mariadb.sql.gz
     
       -db_engine string
         	db engine (sqlite or mariadb) (default "sqlite")
       -db_name string
         	db name
       -mariadb_host string
         	 (default "localhost")
       -mariadb_port string
         	 (default "3306")
       -mariadb_protocol string
         	 (default "tcp")
       -mariadb_user string
         	 (default "speechoid")
       -sqlite_folder string
         	
     
