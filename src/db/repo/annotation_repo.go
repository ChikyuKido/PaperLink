package repo

import (
	"paperlink/db/entity"
)

type AnnotationRepo struct {
	*Repository[entity.Annotation]
}

func newAnnotationRepo() *AnnotationRepo {
	return &AnnotationRepo{NewRepository[entity.Annotation]()}
}

var Annotation = newAnnotationRepo()
