package query_test

import (
	"net/url"
	"testing"

	"github.com/securehaven/query"
	"github.com/stretchr/testify/assert"
)

type exampleUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type examplePost struct {
	Id     int         `json:"id"`
	Title  string      `json:"title"`
	Author exampleUser `json:"author"`
}

var parser = query.MustParser(query.NewParser[examplePost]())

func TestReadme(t *testing.T) {
	queryValues, _ := url.ParseQuery("limit=10&offset=0&sort=id:asc&select=title,author.first_name,author.last_name&id=gt:1")
	q, err := parser.Parse(queryValues)

	assert.NoError(t, err, "should not return an error")

	// fmt.Printf("%+v", q)

	_ = q

	// filtered := q.Filter(Data{
	// 	Id:        3,
	// 	FirstName: "John",
	// 	LastName:  "Doe",
	// })

	// expected := map[string]any{
	// 	"first_name": query.NewNull("John", true),
	// 	"last_name":  query.NewNull("Doe", true),
	// }

	// if !reflect.DeepEqual(filtered, expected) {
	// 	t.Errorf("unexpected filtered result: %v", filtered)
	// }
}

func TestLimit(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("valid", func(t *testing.T) {
		values.Set(query.ParamLimit, "30")

		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, 30, q.Limit, "limit should be equal")
	})

	t.Run("fallback-min", func(t *testing.T) {
		values.Set(query.ParamLimit, "-10")

		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, query.DefaultBaseLimit, q.Limit, "limit should be equal to DefaultBaseLimit")
	})

	t.Run("fallback-max", func(t *testing.T) {
		values.Set(query.ParamLimit, "100000000")

		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, query.DefaultMaxLimit, q.Limit, "limit should be equal to DefaultMaxLimit")
	})
}

func TestOffset(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("valid", func(t *testing.T) {
		values.Set(query.ParamOffset, "8")

		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, 8, q.Offset, "offset should be equal")
	})

	t.Run("fallback", func(t *testing.T) {
		values.Set(query.ParamOffset, "-1")

		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, query.DefaultBaseOffset, q.Offset, "offset should be equal to DefaultBaseOffset")
	})
}

func TestSelect(t *testing.T) {
	values := make(url.Values, 0)

	t.Run("single", func(t *testing.T) {
		values.Set(query.ParamSelect, "id")

		expected := []string{"id"}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Select, "selected fields should be equal")
	})

	t.Run("multiple", func(t *testing.T) {
		values.Set(query.ParamSelect, "id,title,author.first_name,author.last_name")

		expected := []string{"id", "title", "author.first_name", "author.last_name"}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Select, "selected fields should be equal")
	})

	t.Run("unknown-field", func(t *testing.T) {
		values.Set(query.ParamSelect, "id,title,created_at")

		expected := []string{"id", "title"}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Select, "selected fields should be equal")
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

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Sortings, "sortings should be equal")
	})

	t.Run("single", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc")

		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
		}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Sortings, "sortings should be equal")
	})

	t.Run("multiple", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc,author.first_name:asc_nulls_first")

		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
			{Field: "author.first_name", Order: query.OrderAscNullsFirst},
		}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Sortings, "sortings should be equal")
	})

	t.Run("unknown-field", func(t *testing.T) {
		values.Set(query.ParamSort, "id:desc,author.first_name:asc_nulls_first,created_at:desc_nulls_last")

		expected := []query.Sorting{
			{Field: "id", Order: query.OrderDesc},
			{Field: "author.first_name", Order: query.OrderAscNullsFirst},
		}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Sortings, "sortings should be equal")
	})
}

func TestFilter(t *testing.T) {
	t.Run("field-only", func(t *testing.T) {
		values := make(url.Values, 0)
		values.Set("author.first_name", "Joe")

		expected := []query.Filtering{
			{Field: "author.first_name", Filter: query.FilterEquals, Value: "Joe"},
		}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Filterings, "filterings should be equal")
	})

	t.Run("single", func(t *testing.T) {
		values := make(url.Values, 0)
		values.Set("author.first_name", "neq:Joe")

		expected := []query.Filtering{
			{Field: "author.first_name", Filter: query.FilterNotEquals, Value: "Joe"},
		}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Filterings, "filterings should be equal")
	})

	t.Run("multiple", func(t *testing.T) {
		values := make(url.Values, 0)
		values.Set("id", "5")
		values.Set("author.first_name", "neq:Joe")

		expected := []query.Filtering{
			{Field: "id", Filter: query.FilterEquals, Value: 5},
			{Field: "author.first_name", Filter: query.FilterNotEquals, Value: "Joe"},
		}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Filterings, "filterings should be equal")
	})

	t.Run("like", func(t *testing.T) {
		values, _ := url.ParseQuery("author.first_name=like:Jo%25")
		expected := []query.Filtering{
			{Field: "author.first_name", Filter: query.FilterLike, Value: "Jo%"},
		}
		q, err := parser.Parse(values)

		assert.NoError(t, err, "should not return an error")
		assert.Equal(t, expected, q.Filterings, "filterings should be equal")
	})
}
