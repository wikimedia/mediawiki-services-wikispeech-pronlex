package dbapi

// SchemaVersion defines the version of the schema structure. It is used for validating databases against the current version number. It will be updated manually when the structure of the schema/database is changed. Versions with the same prefix (e.g., 3 and 3.1) are compatible.
const SchemaVersion = "3.1"
