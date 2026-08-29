package postgres

// SQLState is a PostgreSQL error code value.
type SQLState string

// PostgreSQL SQLSTATE error codes used by persistence adapters.
const (
	UniqueViolation     SQLState = "23505"
	ForeignKeyViolation SQLState = "23503"
	CheckViolation      SQLState = "23514"
	NotNullViolation    SQLState = "23502"
	InvalidTextValue    SQLState = "22P02"
)
