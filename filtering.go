package query

type Filtering struct {
	Field  string
	Filter string
	Value  string
}

const (
	FilterEquals           = "eq"
	FilterLessThan         = "lt"
	FilterGreaterThan      = "gt"
	FilterLessThanEquals   = "lte"
	FilterGreateThanEquals = "gte"
	FilterLike             = "like"
	FilterNotEquals        = "neq"
)

var (
	allowedFilterValues = []string{
		FilterEquals,
		FilterLessThan,
		FilterGreaterThan,
		FilterLessThanEquals,
		FilterGreateThanEquals,
		FilterLike,
		FilterNotEquals,
	}
)
