package query

type Filtering struct {
	Field  string
	Filter string
	Value  any
}

const (
	FilterEquals           = "eq"
	FilterLessThan         = "lt"
	FilterGreaterThan      = "gt"
	FilterLessThanEquals   = "lte"
	FilterGreateThanEquals = "gte"
	FilterLike             = "like"
	FilterNot              = "not"
)

var (
	allowedFilterValues = []string{
		FilterEquals,
		FilterLessThan,
		FilterGreaterThan,
		FilterLessThanEquals,
		FilterGreateThanEquals,
		FilterLike,
		FilterNot,
	}
)
