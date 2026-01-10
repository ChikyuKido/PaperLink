package d4s

import (
	"paperlink/db/entity"
	"paperlink/util"
)

type Book struct {
	Name     string                     `json:"name"`
	DataCode string                     `json:"dataCode"`
	DataId   string                     `json:"dataId"`
	UUID     string                     `json:"-"`
	Account  *entity.Digi4SchoolAccount `json:"-"`
}

var log = util.GroupLog("SERVICE_DIGI4SCHOOL")
