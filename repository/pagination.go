package repository

//Page ...
type Page struct {
	Offset  uint   `json:"offset"`
	Amount  uint   `json:"amount"`
	OrderBy string `json:"orderBy"`
	Sort    string `json:"sort"`
}
