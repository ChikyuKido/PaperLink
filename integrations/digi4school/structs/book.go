package structs

type Book struct {
	Name      string `json:"name"`
	DataCode  string `json:"dataCode"`
	DataId    string `json:"dataId"`
	EbookPlus bool   `json:"ebook_plus"`
	Publisher string `json:"publisher"`
}
