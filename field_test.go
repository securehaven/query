package query

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldsFromStruct(t *testing.T) {
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
