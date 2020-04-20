-- $ sudo mysql -u root < mariadb_setup.sql

CREATE USER IF NOT EXISTS 'speechoid'@'localhost';
-- DROP DATABASE IF EXISTS speechoid_pronlex_test1;

-- lexserver demo db
CREATE DATABASE lexserver_testdb;
GRANT ALL PRIVILEGES ON lexserver_testdb.* TO 'speechoid'@'localhost' ;

-- Test_insertEntries
CREATE DATABASE speechoid_pronlex_test1;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test1.* TO 'speechoid'@'localhost' ;

-- Test_ImportLexiconFile
CREATE DATABASE speechoid_pronlex_test2;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test2.* TO 'speechoid'@'localhost' ;

-- Test_ImportLexiconFileWithDupLines
CREATE DATABASE speechoid_pronlex_test3;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test3.* TO 'speechoid'@'localhost' ;

-- Test_ImportLexiconFileInvalid
CREATE DATABASE speechoid_pronlex_test4;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test4.* TO 'speechoid'@'localhost' ;

-- Test_ImportLexiconFileGz
CREATE DATABASE speechoid_pronlex_test5;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test5.* TO 'speechoid'@'localhost' ;

-- Test_UpdateComments
CREATE DATABASE speechoid_pronlex_test6;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test6.* TO 'speechoid'@'localhost' ;

-- Test_ValidationRuleLike
CREATE DATABASE speechoid_pronlex_test7;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test7.* TO 'speechoid'@'localhost' ;

-- Test_DBManager
CREATE DATABASE speechoid_pronlex_test8;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test8.* TO 'speechoid'@'localhost' ;
CREATE DATABASE speechoid_pronlex_test9;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test9.* TO 'speechoid'@'localhost' ;

-- Test_MoveNewEntries
CREATE DATABASE speechoid_pronlex_test10;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test10.* TO 'speechoid'@'localhost' ;

-- TestEntryTag1
CREATE DATABASE speechoid_pronlex_test11;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test11.* TO 'speechoid'@'localhost' ;

-- TestEntryTag2
CREATE DATABASE speechoid_pronlex_test12;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test12.* TO 'speechoid'@'localhost' ;

-- Test_Validation1
CREATE DATABASE speechoid_pronlex_test13;
GRANT ALL PRIVILEGES ON speechoid_pronlex_test13.* TO 'speechoid'@'localhost' ;
