package query

type Query struct {
	Limit      int
	Offset     int
	Select     []string
	Sortings   []Sorting
	Filterings []Filtering
}
