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

type Parser struct {
	allowedFields []string
	maxLimit      int
	baseLimit     int
	baseOffset    int
}

func NewParser(allowedFields []string) *Parser {
	return &Parser{
		allowedFields: allowedFields,
		maxLimit:      DefaultMaxLimit,
		baseLimit:     DefaultBaseLimit,
		baseOffset:    DefaultBaseOffset,
	}
}

func (p *Parser) Parse(query url.Values) Query {
	return Query{
		Limit:      p.parseLimit(query.Get(ParamLimit)),
		Offset:     p.parseOffset(query.Get(ParamOffset)),
		Select:     p.parseSelect(query.Get(ParamSelect)),
		Sortings:   p.parseSort(query.Get(ParamSort)),
		Filterings: p.parseFilter(query),
	}
}

func (p *Parser) parseFilter(values url.Values) []Filtering {
	filterings := make([]Filtering, 0, len(p.allowedFields))

	for _, field := range p.allowedFields {
		raw := values.Get(field)

		if len(raw) == 0 {
			continue
		}

		parts := p.splitClean(raw, SeparatorFilter, 2)

		if len(parts) == 1 {
			filterings = append(filterings, Filtering{
				Field:  field,
				Filter: FilterEquals,
				Value:  parts[0], // TODO: Parse into correct type
			})
			continue
		}

		if !p.isAllowedFilterValue(parts[0]) {
			continue
		}

		filterings = append(filterings, Filtering{
			Field:  field,
			Filter: parts[0],
			Value:  parts[1],
		})
	}

	return filterings
}

func (p *Parser) parseSort(raw string) []Sorting {
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

func (p *Parser) parseSelect(raw string) []string {
	fields := p.splitClean(raw, SeparatorField, -1)

	return slices.DeleteFunc(fields, func(field string) bool {
		return !p.isAllowedField(field)
	})
}

func (p *Parser) splitClean(raw string, sep string, n int) []string {
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

func (p *Parser) parseLimit(raw string) int {
	limit := p.parseInt(raw, p.baseLimit)

	if limit <= 0 {
		return p.baseLimit
	}

	if limit > p.maxLimit {
		return p.maxLimit
	}

	return limit
}

func (p *Parser) parseOffset(raw string) int {
	offset := p.parseInt(raw, p.baseOffset)

	if offset < 0 {
		return p.baseOffset
	}

	return offset
}

func (p *Parser) parseInt(raw string, fallback int) int {
	if len(raw) == 0 {
		return fallback
	}

	value, err := strconv.Atoi(raw)

	if err != nil {
		return fallback
	}

	return value
}

func (p *Parser) isAllowedField(field string) bool {
	return slices.Contains(p.allowedFields, field)
}

func (p *Parser) isAllowedOrderValue(order string) bool {
	return slices.Contains(allowedOrderValues, order)
}

func (p *Parser) isAllowedFilterValue(filter string) bool {
	return slices.Contains(allowedFilterValues, filter)
}

func (p *Parser) WithMaxLimit(max int) *Parser {
	p.maxLimit = max

	return p
}

func (p *Parser) WithBaseLimit(base int) *Parser {
	p.baseLimit = base

	return p
}

func (p *Parser) WithBaseOffset(base int) *Parser {
	p.baseOffset = base

	return p
}
