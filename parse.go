package query

import (
	"strconv"
	"time"
)

func ParseString(v string) (any, error) {
	return v, nil
}

func ParseBool(v string) (any, error) {
	return strconv.ParseBool(v)
}

func ParseInt(base int, bitSize int) ParseFunc {
	return func(v string) (any, error) {
		r, err := strconv.ParseInt(v, base, bitSize)

		switch bitSize {
		case 0:
			return int(r), err
		case 8:
			return int8(r), err
		case 16:
			return int16(r), err
		case 32:
			return int32(r), err
		default:
			return r, err
		}
	}
}

func ParseUint(base int, bitSize int) ParseFunc {
	return func(v string) (any, error) {
		r, err := strconv.ParseUint(v, base, bitSize)

		switch bitSize {
		case 0:
			return uint(r), err
		case 8:
			return uint8(r), err
		case 16:
			return uint16(r), err
		case 32:
			return uint32(r), err
		default:
			return r, err
		}
	}
}

func ParseFloat(bitSize int) ParseFunc {
	return func(v string) (any, error) {
		r, err := strconv.ParseFloat(v, bitSize)

		switch bitSize {
		case 32:
			return float32(r), err
		default:
			return r, err
		}
	}
}

func ParseTime(layout string) ParseFunc {
	return func(v string) (any, error) {
		return time.Parse(layout, v)
	}
}
