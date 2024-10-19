package query

import "slices"

type Query[S any] struct {
	Limit      int
	Offset     int
	Select     []string
	Sortings   []Sorting
	Filterings []Filtering

	fields map[string]ValueFunc[S]
}

func (q Query[S]) Filter(data S) map[string]any {
	filteredFields := make(map[string]any, len(q.Select))

	for name, valueFunc := range q.fields {
		if len(q.Select) > 0 && !slices.Contains(q.Select, name) {
			continue
		}

		filteredFields[name] = valueFunc(data)
	}

	return filteredFields
}
