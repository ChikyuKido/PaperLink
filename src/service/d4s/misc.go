package d4s

import "paperlink/util"

type Book struct {
	Name     string `json:"name"`
	DataCode string `json:"dataCode"`
	DataId   string `json:"dataId"`
}

var log = util.GroupLog("SERVICE_DIGI4SCHOOL")
