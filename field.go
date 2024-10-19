package query

type ParseFunc func(v string) (any, error)

type ValueFunc[S any] func(S) any

type Field[S any] struct {
	name      string
	parseFunc ParseFunc
	valueFunc ValueFunc[S]
}

func NewField[S any](name string, p ParseFunc, v ValueFunc[S]) Field[S] {
	return Field[S]{
		name:      name,
		parseFunc: p,
		valueFunc: v,
	}
}
