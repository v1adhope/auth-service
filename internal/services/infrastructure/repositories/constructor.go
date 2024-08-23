package repositories

import "github.com/v1adhope/auth-service/pkg/postgresql"

type Repos struct {
	*postgresql.Postgres
}

func New(driver *postgresql.Postgres) *Repos {
	return &Repos{driver}
}
