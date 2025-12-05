package repo

import (
	"paperlink/db/entity"
)

type DocumentUserRepo struct {
	*Repository[entity.DocumentUser]
}

func newDocumentUserRepo() *DocumentUserRepo {
	return &DocumentUserRepo{NewRepository[entity.DocumentUser]()}
}

var DocumentUser = newDocumentUserRepo()
