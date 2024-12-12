package query

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldsFromStruct(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		type example struct {
			Id        int
			Age       int    `json:"-"`
			UserName  string `json:""`
			FirstName string `json:"firstName"`
			LastName  string `json:"last_name"`
		}

		fields, err := getFieldsFromStruct[example]()

		if !assert.NoError(t, err, "unexpected error") {
			assert.ErrorIs(t, err, ErrNoStruct, "generic type must be a struct")
		}

		names := make([]string, len(fields))

		for i, field := range fields {
			names[i] = field.name
		}

		expected := []string{"Id", "UserName", "firstName", "last_name"}

		assert.ElementsMatch(t, expected, names, "unexpected result")
	})

	t.Run("complex", func(t *testing.T) {
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

		fields, err := getFieldsFromStruct[examplePost]()

		if !assert.NoError(t, err, "unexpected error") {
			assert.ErrorIs(t, err, ErrNoStruct, "generic type must be a struct")
		}

		names := make([]string, len(fields))

		for i, field := range fields {
			names[i] = field.name
		}

		expected := []string{"id", "title", "author", "author.id", "author.first_name", "author.last_name"}

		assert.ElementsMatch(t, expected, names, "unexpected result")
	})
}

func TestGetFieldNameFromStructField(t *testing.T) {
	t.Run("tagless", func(t *testing.T) {
		type example struct {
			UserName string
		}

		var e example

		structType := reflect.TypeOf(e)
		structField, exists := structType.FieldByName("UserName")

		assert.Equal(t, true, exists, "field UserName should exist in example struct")

		name := getFieldNameFromStructField(structField)

		assert.Equal(t, "UserName", name)
	})

	t.Run("hidden", func(t *testing.T) {
		type example struct {
			UserName string `json:"-"`
		}

		var e example

		structType := reflect.TypeOf(e)
		structField, exists := structType.FieldByName("UserName")

		assert.Equal(t, true, exists, "field UserName should exist in example struct")

		name := getFieldNameFromStructField(structField)

		assert.Empty(t, name, "field name should be empty")
	})

	t.Run("empty-tag", func(t *testing.T) {
		type example struct {
			UserName string `json:""`
		}

		var e example

		structType := reflect.TypeOf(e)
		structField, exists := structType.FieldByName("UserName")

		assert.Equal(t, true, exists, "field UserName should exist in example struct")

		name := getFieldNameFromStructField(structField)

		assert.Equal(t, "UserName", name)
	})

	t.Run("replace", func(t *testing.T) {
		type example struct {
			UserName string `json:"user_name"`
		}

		var e example

		structType := reflect.TypeOf(e)
		structField, exists := structType.FieldByName("UserName")

		assert.Equal(t, true, exists, "field UserName should exist in example struct")

		name := getFieldNameFromStructField(structField)

		assert.Equal(t, "user_name", name)
	})
}
