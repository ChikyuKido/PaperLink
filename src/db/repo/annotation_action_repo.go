package repo

import (
	"paperlink/db/entity"
)

type AnnotationActionRepo struct {
	*Repository[entity.AnnotationAction]
}

func newAnnotationActionRepo() *AnnotationActionRepo {
	return &AnnotationActionRepo{NewRepository[entity.AnnotationAction]()}
}

var AnnotationAction = newAnnotationActionRepo()
