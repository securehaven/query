package query

import (
	"net/url"
	"slices"
	"strconv"
	"strings"
)

const (
	DefaultMaxLimit   = 100
	DefaultBaseLimit  = 10
	DefaultBaseOffset = 0
)

var (
	ParamLimit  = "limit"
	ParamOffset = "offset"
	ParamSelect = "select"
	ParamSort   = "sort"

	SeparatorField  = ","
	SeparatorFilter = ":"
)

type Parser[S any] struct {
	fields     []Field[S]
	maxLimit   int
	baseLimit  int
	baseOffset int
}

func NewParser[S any](fields ...Field[S]) *Parser[S] {
	return &Parser[S]{
		fields:     fields,
		maxLimit:   DefaultMaxLimit,
		baseLimit:  DefaultBaseLimit,
		baseOffset: DefaultBaseOffset,
	}
}

func (p *Parser[S]) Parse(v url.Values) (Query[S], error) {
	filterings, err := p.parseFilter(v)

	return Query[S]{
		Limit:      p.parseLimit(v.Get(ParamLimit)),
		Offset:     p.parseOffset(v.Get(ParamOffset)),
		Select:     p.parseSelect(v.Get(ParamSelect)),
		Sortings:   p.parseSort(v.Get(ParamSort)),
		Filterings: filterings,
		fields:     p.getQueryFields(),
	}, err
}

func (p *Parser[S]) getQueryFields() map[string]ValueFunc[S] {
	queryFields := make(map[string]ValueFunc[S], len(p.fields))

	for _, field := range p.fields {
		queryFields[field.name] = field.valueFunc
	}

	return queryFields
}

func (p *Parser[S]) parseFilter(values url.Values) ([]Filtering, error) {
	filterings := make([]Filtering, 0, len(p.fields))

	for _, field := range p.fields {
		raw := values.Get(field.name)

		if len(raw) == 0 {
			continue
		}

		parts := p.splitClean(raw, SeparatorFilter, 2)
		filtering, err := p.newFiltering(field.name, field.parseFunc, parts)

		if err != nil {
			return filterings, err
		}

		filterings = append(filterings, filtering)
	}

	return filterings, nil
}

func (p *Parser[S]) newFiltering(field string, parse ParseFunc, parts []string) (Filtering, error) {
	var err error

	filtering := Filtering{
		Field: field,
	}

	if len(parts) == 1 {
		filtering.Filter = FilterEquals
		filtering.Value, err = parse(parts[0])
	} else {
		if !p.isAllowedFilterValue(parts[0]) {
			return Filtering{}, nil
		}

		filtering.Filter = parts[0]
		filtering.Value, err = parse(parts[1])
	}

	return filtering, err
}

func (p *Parser[S]) parseSort(raw string) []Sorting {
	rawParts := p.splitClean(raw, SeparatorField, -1)
	sortings := make([]Sorting, 0, len(rawParts))

	for _, rawPart := range rawParts {
		parts := p.splitClean(rawPart, SeparatorFilter, 2)

		if len(parts) < 1 {
			continue
		}

		if !p.isAllowedField(parts[0]) {
			continue
		}

		if len(parts) < 2 {
			sortings = append(sortings, Sorting{
				Field: parts[0],
				Order: OrderAsc,
			})
			continue
		}

		if !p.isAllowedOrderValue(parts[1]) {
			continue
		}

		sortings = append(sortings, Sorting{
			Field: parts[0],
			Order: parts[1],
		})
	}

	return sortings
}

func (p *Parser[S]) parseSelect(raw string) []string {
	fields := p.splitClean(raw, SeparatorField, -1)

	return slices.DeleteFunc(fields, func(field string) bool {
		return !p.isAllowedField(field)
	})
}

func (p *Parser[S]) parseLimit(raw string) int {
	limit := p.parseInt(raw, p.baseLimit)

	if limit <= 0 {
		return p.baseLimit
	}

	if limit > p.maxLimit {
		return p.maxLimit
	}

	return limit
}

func (p *Parser[S]) parseOffset(raw string) int {
	offset := p.parseInt(raw, p.baseOffset)

	if offset < 0 {
		return p.baseOffset
	}

	return offset
}

func (p *Parser[S]) parseInt(raw string, fallback int) int {
	if len(raw) == 0 {
		return fallback
	}

	value, err := strconv.Atoi(raw)

	if err != nil {
		return fallback
	}

	return value
}

func (p *Parser[S]) splitClean(raw string, sep string, n int) []string {
	rawParts := strings.SplitN(raw, sep, n)
	parts := make([]string, 0, len(rawParts))

	for _, rawPart := range rawParts {
		part := strings.TrimSpace(rawPart)

		if len(part) > 0 {
			parts = append(parts, part)
		}
	}

	return parts
}

func (p *Parser[S]) isAllowedField(field string) bool {
	return slices.ContainsFunc(p.fields, func(f Field[S]) bool {
		return f.name == field
	})
}

func (p *Parser[S]) isAllowedOrderValue(order string) bool {
	return slices.Contains(allowedOrderValues, order)
}

func (p *Parser[S]) isAllowedFilterValue(filter string) bool {
	return slices.Contains(allowedFilterValues, filter)
}

func (p *Parser[S]) WithMaxLimit(max int) *Parser[S] {
	p.maxLimit = max

	return p
}

func (p *Parser[S]) WithBaseLimit(base int) *Parser[S] {
	p.baseLimit = base

	return p
}

func (p *Parser[S]) WithBaseOffset(base int) *Parser[S] {
	p.baseOffset = base

	return p
}
