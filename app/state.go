package app

type EntryMeta struct {
	Title      string `json:"title"`
	Date       string `json:"date"`
	CreatedAt  string `json:"created_at"`
	LastEdited string `json:"last_edited"`
	Path       string `json:"path"`
}
