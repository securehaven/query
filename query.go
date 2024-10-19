package query

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

	for _, name := range q.Select {
		valueFunc, ok := q.fields[name]

		if ok {
			filteredFields[name] = valueFunc(data)
		}
	}

	return filteredFields
}
