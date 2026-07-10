package types

import (
	"time"
	"database/sql"
)


type FileMeta struct{
	ID string
	Name string
	Hash string //hex encoded SHA-256
	Owner string
	CreatedAt time.Time
	Size int64
}

type MySQLMetadata struct{
	db *sql.DB
}