package query

import (
	"net/url"
	"strconv"
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
	return newBuilder(p.allowedFields).
		withLimit(p.parseLimit(query.Get(ParamLimit))).
		withOffset(p.parseOffset(query.Get(ParamOffset))).
		get()
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
