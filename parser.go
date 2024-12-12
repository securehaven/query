package query

import (
	"log"
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

	SeparatorSelector = "."
	SeparatorField    = ","
	SeparatorFilter   = ":"
)

type Parser struct {
	fields     []field
	maxLimit   int
	baseLimit  int
	baseOffset int
}

func MustParser(p *Parser, err error) *Parser {
	if err != nil {
		log.Fatalf("could not create query parser: %v", err)
	}

	return p
}

func NewParser[T any]() (*Parser, error) {
	fields, err := getFieldsFromStruct[T]()

	return &Parser{
		fields:     fields,
		maxLimit:   DefaultMaxLimit,
		baseLimit:  DefaultBaseLimit,
		baseOffset: DefaultBaseOffset,
	}, err
}

func (p *Parser) Parse(v url.Values, excludes ...string) (Query, error) {
	filterings, err := p.parseFilter(v)

	return Query{
		Limit:      p.parseLimit(v.Get(ParamLimit)),
		Offset:     p.parseOffset(v.Get(ParamOffset)),
		Select:     p.parseSelect(v.Get(ParamSelect)),
		Sortings:   p.parseSort(v.Get(ParamSort)),
		Filterings: filterings,
	}, err
}

func (p *Parser) parseFilter(values url.Values) ([]Filtering, error) {
	filterings := make([]Filtering, 0, len(p.fields))
	parsingError := ParsingError{}

	for _, field := range p.fields {
		raw := values.Get(field.name)

		if len(raw) == 0 {
			continue
		}

		parts := p.splitClean(raw, SeparatorFilter, 2)
		filtering, err := p.newFiltering(field.name, field.parseFunc, parts)

		if err != nil {
			parsingError.Errors = append(parsingError.Errors, err)
		}

		filterings = append(filterings, filtering)
	}

	if len(parsingError.Errors) > 0 {
		return filterings, parsingError
	}

	return filterings, nil
}

func (p *Parser) newFiltering(field string, parse ParseFunc, parts []string) (Filtering, error) {
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

func (p *Parser) isAllowedField(name string) bool {
	return slices.ContainsFunc(p.fields, func(f field) bool {
		return f.name == name
	})
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
