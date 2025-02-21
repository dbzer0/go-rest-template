package models

// Paging представляет собой модель для пагинации в API
type Paging struct {
	Skip      *int64
	Limit     *int64
	SortKey   string
	SortVal   int
	Condition interface{}
}
