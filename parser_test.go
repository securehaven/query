package query_test

import (
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/securehaven/query"
)

type Data struct {
	Id        int
	FirstName string
	LastName  string
}

var parser = query.NewParser(
	query.NewField("id", query.ParseInt(0, 0), func(d Data) any {
		return query.NewNull(d.Id, d.Id > 0)
	}),
	query.NewField("first_name", query.ParseString, func(d Data) any {
		return query.NewNull(d.FirstName, len(d.FirstName) > 0)
	}),
	query.NewField("last_name", query.ParseString, func(d Data) any {
		return query.NewNull(d.LastName, len(d.LastName) > 0)
	}),
)

func TestReadme(t *testing.T) {
	queryValues, _ := url.ParseQuery("limit=10&offset=0&sort=id:asc&select=first_name,last_name&id=gt:1")
	q, err := parser.Parse(queryValues)

	if err != nil {
		t.Errorf("failed to parse some value: %v", err)
	}

	filtered := q.Filter(Data{
		Id:        3,
		FirstName: "John",
		LastName:  "Doe",
	})

	expected := map[string]any{
		"first_name": query.NewNull("John", true),
		"last_name":  query.NewNull("Doe", true),
	}

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("unexpected filtered result: %v", filtered)
	}
}

func TestLimit(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("valid", func(t *testing.T) {
		values.Set(query.ParamLimit, "30")

		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if q.Limit != 30 {
			t.Errorf("%q received, expected limit to be %q", q.Limit, 30)
		}
	})

	t.Run("fallback-min", func(t *testing.T) {
		values.Set(query.ParamLimit, "-10")

		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if q.Limit != query.DefaultBaseLimit {
			t.Errorf("%q received, expected limit to be %q", q.Limit, query.DefaultBaseLimit)
		}
	})

	t.Run("fallback-max", func(t *testing.T) {
		values.Set(query.ParamLimit, "100000000")

		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if q.Limit != query.DefaultMaxLimit {
			t.Errorf("%q received, expected limit to be %q", q.Limit, query.DefaultMaxLimit)
		}
	})
}

func TestOffset(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("valid", func(t *testing.T) {
		values.Set(query.ParamOffset, "8")

		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if q.Offset != 8 {
			t.Errorf("%q received, expected offset to be %q", q.Offset, 8)
		}
	})

	t.Run("fallback", func(t *testing.T) {
		values.Set(query.ParamOffset, "-1")

		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if q.Offset != query.DefaultBaseOffset {
			t.Errorf("%q received, expected offset to be %q", q.Offset, query.DefaultBaseOffset)
		}
	})
}

func TestSelect(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("single", func(t *testing.T) {
		values.Set(query.ParamSelect, "id")

		expected := []string{"id"}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Select, expected) {
			t.Errorf("%q received, expected select to be %q",
				strings.Join(q.Select, ","),
				strings.Join(expected, ","),
			)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		values.Set(query.ParamSelect, "id,first_name,last_name")

		expected := []string{"id", "first_name", "last_name"}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Select, expected) {
			t.Errorf("%q received, expected select to be %q",
				strings.Join(q.Select, ","),
				strings.Join(expected, ","),
			)
		}
	})

	t.Run("unknown-field", func(t *testing.T) {
		values.Set(query.ParamSelect, "id,first_name,last_name,created_at")

		expected := []string{"id", "first_name", "last_name"}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

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

		expected := []query.Sorting{
			{Field: "id", Order: query.OrderAsc},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})

	t.Run("single", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc")

		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc,first_name:asc_nulls_first")

		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
			{Field: "first_name", Order: query.OrderAscNullsFirst},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})

	t.Run("unknown-field", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc,first_name:asc_nulls_first,created_at:desc_nulls_last")

		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
			{Field: "first_name", Order: query.OrderAscNullsFirst},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Sortings, expected) {
			t.Errorf("unexpected sorting result: %+v", q.Sortings)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("field-only", func(t *testing.T) {
		values := make(url.Values, 0)
		values.Set("first_name", "Joe")

		expected := []query.Filtering{
			{Field: "first_name", Filter: query.FilterEquals, Value: "Joe"},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Filterings, expected) {
			t.Errorf("unexpected filtering result: %+v", q.Filterings)
		}
	})

	t.Run("single", func(t *testing.T) {
		values := make(url.Values, 0)
		values.Set("first_name", "neq:Joe")

		expected := []query.Filtering{
			{Field: "first_name", Filter: query.FilterNotEquals, Value: "Joe"},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Filterings, expected) {
			t.Errorf("unexpected filtering result: %+v", q.Filterings)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		values := make(url.Values, 0)
		values.Set("id", "5")
		values.Set("first_name", "neq:Joe")

		expected := []query.Filtering{
			{Field: "id", Filter: query.FilterEquals, Value: 5},
			{Field: "first_name", Filter: query.FilterNotEquals, Value: "Joe"},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Filterings, expected) {
			t.Errorf("unexpected filtering result: %+v", q.Filterings)
		}
	})

	t.Run("like", func(t *testing.T) {
		values, _ := url.ParseQuery("first_name=like:Jo%25")
		expected := []query.Filtering{
			{Field: "first_name", Filter: query.FilterLike, Value: "Jo%"},
		}
		q, err := parser.Parse(values)

		if err != nil {
			t.Error(err)
		}

		if !slices.Equal(q.Filterings, expected) {
			t.Errorf("unexpected filtering result: %+v", q.Filterings)
		}
	})
}
