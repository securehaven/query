package query_test

import (
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/securehaven/query"
)

var parser = query.NewParser([]string{"id", "first_name", "last_name"})

func TestLimit(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("valid", func(t *testing.T) {
		values.Set(query.ParamLimit, "30")

		q := parser.Parse(values)

		if q.Limit != 30 {
			t.Errorf("%q received, expected limit to be %q", q.Limit, 30)
		}
	})

	t.Run("fallback-min", func(t *testing.T) {
		values.Set(query.ParamLimit, "-10")

		q := parser.Parse(values)

		if q.Limit != query.DefaultBaseLimit {
			t.Errorf("%q received, expected limit to be %q", q.Limit, query.DefaultBaseLimit)
		}
	})

	t.Run("fallback-max", func(t *testing.T) {
		values.Set(query.ParamLimit, "100000000")

		q := parser.Parse(values)

		if q.Limit != query.DefaultMaxLimit {
			t.Errorf("%q received, expected limit to be %q", q.Limit, query.DefaultMaxLimit)
		}
	})
}

func TestOffset(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("valid", func(t *testing.T) {
		values.Set(query.ParamOffset, "8")

		q := parser.Parse(values)

		if q.Offset != 8 {
			t.Errorf("%q received, expected offset to be %q", q.Offset, 8)
		}
	})

	t.Run("fallback", func(t *testing.T) {
		values.Set(query.ParamOffset, "-1")

		q := parser.Parse(values)

		if q.Offset != query.DefaultBaseOffset {
			t.Errorf("%q received, expected offset to be %q", q.Offset, query.DefaultBaseOffset)
		}
	})
}

func TestSelect(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("single", func(t *testing.T) {
		values.Set(query.ParamSelect, "id")

		q := parser.Parse(values)
		expected := []string{"id"}

		if !slices.Equal(q.Select, expected) {
			t.Errorf("%q received, expected select to be %q",
				strings.Join(q.Select, ","),
				strings.Join(expected, ","),
			)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		values.Set(query.ParamSelect, "id,first_name,last_name")

		q := parser.Parse(values)
		expected := []string{"id", "first_name", "last_name"}

		if !slices.Equal(q.Select, expected) {
			t.Errorf("%q received, expected select to be %q",
				strings.Join(q.Select, ","),
				strings.Join(expected, ","),
			)
		}
	})

	t.Run("unknown-field", func(t *testing.T) {
		values.Set(query.ParamSelect, "id,first_name,last_name,created_at")

		q := parser.Parse(values)
		expected := []string{"id", "first_name", "last_name"}

		if !slices.Equal(q.Select, expected) {
			t.Errorf("%q received, expected select to be %q",
				strings.Join(q.Select, ","),
				strings.Join(expected, ","),
			)
		}
	})
}

func TestSort(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("field-only", func(t *testing.T) {
		values.Set(query.ParamSort, "id")

		q := parser.Parse(values)
		expected := []query.Sorting{
			{Field: "id", Order: query.OrderAsc},
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})

	t.Run("single", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc")

		q := parser.Parse(values)
		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc,first_name:asc_nulls_first")

		q := parser.Parse(values)
		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
			{Field: "first_name", Order: query.OrderAscNullsFirst},
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})

	t.Run("unknown-field", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc,first_name:asc_nulls_first,created_at:desc_nulls_last")

		q := parser.Parse(values)
		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
			{Field: "first_name", Order: query.OrderAscNullsFirst},
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})
}

func TestFilter(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("field-only", func(t *testing.T) {
		values.Set("first_name", "Joe")

		q := parser.Parse(values)
		expected := []query.Filtering{
			{Field: "first_name", Filter: query.FilterEquals, Value: "Joe"},
		}

		if !slices.Equal(q.Filterings, expected) {
			t.Errorf("unexpected filtering result: %+v", q.Filterings)
		}
	})

	t.Run("single", func(t *testing.T) {
		values.Set("first_name", "neq:Joe")

		q := parser.Parse(values)
		expected := []query.Filtering{
			{Field: "first_name", Filter: query.FilterNotEquals, Value: "Joe"},
		}

		if !slices.Equal(q.Filterings, expected) {
			t.Errorf("unexpected filtering result: %+v", q.Filterings)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		values.Set("id", "5")
		values.Set("first_name", "neq:Joe")

		q := parser.Parse(values)
		expected := []query.Filtering{
			{Field: "id", Filter: query.FilterEquals, Value: "5"},
			{Field: "first_name", Filter: query.FilterNotEquals, Value: "Joe"},
		}

		if !slices.Equal(q.Filterings, expected) {
			t.Errorf("unexpected filtering result: %+v", q.Filterings)
		}
	})
}
