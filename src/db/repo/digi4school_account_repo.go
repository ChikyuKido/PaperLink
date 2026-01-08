package repo

import (
	"paperlink/db/entity"
)

type Digi4SchoolAccountRepo struct {
	*Repository[entity.Digi4SchoolAccount]
}

func newDigi4SchoolAccountRepo() *Digi4SchoolAccountRepo {
	return &Digi4SchoolAccountRepo{NewRepository[entity.Digi4SchoolAccount]()}
}

var Digi4SchoolAccount = newDigi4SchoolAccountRepo()
