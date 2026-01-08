package repo

import (
	"paperlink/db/entity"
)

type Digi4SchoolBookRepo struct {
	*Repository[entity.Digi4SchoolBook]
}

func newDigi4SchoolBookRepo() *Digi4SchoolBookRepo {
	return &Digi4SchoolBookRepo{NewRepository[entity.Digi4SchoolBook]()}
}

var Digi4SchoolBook = newDigi4SchoolBookRepo()
