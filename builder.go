package query

import "slices"

type builder struct {
	query         Query
	allowedFields []string
}

func NewBuilder(allowedFields []string) *builder {
	return &builder{
		allowedFields: allowedFields,
	}
}

func (b *builder) WithLimit(limit int, base int, max int) *builder {
	if limit <= 0 {
		b.query.Limit = base
	} else if limit > max {
		b.query.Limit = max
	} else {
		b.query.Limit = limit
	}

	return b
}

func (b *builder) WithOffset(offset int, base int) *builder {
	if offset < 0 {
		b.query.Offset = base
	} else {
		b.query.Offset = offset
	}

	return b
}

func (b *builder) WithSelect(fields []string) *builder {
	for _, field := range fields {
		if b.isAllowedField(field) {
			b.query.Select = append(b.query.Select, field)
		}
	}

	return b
}

func (b *builder) WithSorting(sortings []Sorting) *builder {
	for _, sorting := range sortings {
		if !b.isAllowedField(sorting.Field) {
			continue
		}

		if !b.isAllowedOrderValue(sorting.Order) {
			continue
		}

		b.query.Sortings = append(b.query.Sortings, sorting)
	}

	return b
}

func (b *builder) WithFiltering(filterings []Filtering) *builder {
	for _, filtering := range filterings {
		if !b.isAllowedField(filtering.Field) {
			continue
		}

		if !b.isAllowedFilterValue(filtering.Filter) {
			continue
		}

		b.query.Filterings = append(b.query.Filterings, filtering)
	}

	return b
}

func (b *builder) isAllowedField(field string) bool {
	return slices.Contains(b.allowedFields, field)
}

func (b *builder) isAllowedOrderValue(order string) bool {
	return slices.Contains(allowedOrderValues, order)
}

func (b *builder) isAllowedFilterValue(filter string) bool {
	return slices.Contains(allowedFilterValues, filter)
}
