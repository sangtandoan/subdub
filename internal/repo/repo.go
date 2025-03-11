package repo

import (
	"database/sql"
	"time"
)

const QueryTimeOut = time.Second * 10

type Repo struct {
	User UserRepo
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		User: NewUserRepo(db),
	}
}
