package query

type Sorting struct {
	Field string
	Order string
}

const (
	OrderAsc           = "asc"
	OrderAscNullsFirst = "asc_nulls_first"
	OrderAscNullsLast  = "desc_nulls_last"

	OrderDesc           = "desc"
	OrderDescNullsFirst = "desc_nulls_first"
	OrderDescNullsLast  = "desc_nulls_last"
)

var (
	allowedOrderValues = []string{
		OrderAsc,
		OrderAscNullsFirst,
		OrderAscNullsLast,
		OrderDesc,
		OrderDescNullsFirst,
		OrderDescNullsLast,
	}
)
