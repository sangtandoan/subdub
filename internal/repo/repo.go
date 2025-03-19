package repo

import (
	"database/sql"
	"time"
)

const QueryTimeOut = time.Second * 10

type Repo struct {
	User         UserRepo
	Subscription SubscriptionRepo
	Session      SessionRepo
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		User:         NewUserRepo(db),
		Subscription: NewSubsciptionRepo(db),
		Session:      NewSessionRepo(db),
	}
}
