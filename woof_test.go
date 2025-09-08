package gowoof

import (
	"testing"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/go-viper/mapstructure/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTable_ReturnsExpectedData(t *testing.T) {
	t.Parallel()
	// Arrange
	type TestType struct {
		String  string
		Number  int
		Boolean bool
	}

	table := &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{{Value: "string"}, {Value: "number"}, {Value: "boolean"}}},
			{Cells: []*messages.PickleTableCell{{Value: "abc"}, {Value: "123"}, {Value: "true"}}},
			{Cells: []*messages.PickleTableCell{{Value: "def"}, {Value: "456"}, {Value: "false"}}},
		},
	}

	// Act
	actual, err := ParseTable[TestType](table)

	// Assert
	require.NoError(t, err)

	expected := []TestType{
		{String: "abc", Number: 123, Boolean: true},
		{String: "def", Number: 456, Boolean: false},
	}

	assert.Equal(t, expected, actual)
}

func TestParseTable_ReturnsOptionError(t *testing.T) {
	t.Parallel()
	// Arrange
	type TestType struct{}

	customOption := func(*Config) error {
		return assert.AnError
	}

	// Act
	actual, err := ParseTable[TestType](&godog.Table{}, customOption)

	// Assert
	require.ErrorIs(t, err, assert.AnError)
	assert.Empty(t, actual)
}

func TestParseTable_ReturnsErrorIfGenericIsNotAStruct(t *testing.T) {
	t.Parallel()
	// Act
	actual, err := ParseTable[string](nil)

	// Assert
	require.ErrorIs(t, err, ErrInvalidInput)
	require.ErrorContains(t, err, "generic type is not a struct")
	assert.Empty(t, actual)
}

func TestParseTable_ReturnsErrorOnNoRows(t *testing.T) {
	t.Parallel()
	// Arrange
	type TestType struct{}

	table := &godog.Table{
		Rows: []*messages.PickleTableRow{},
	}

	// Act
	actual, err := ParseTable[TestType](table)

	// Assert
	require.ErrorIs(t, err, ErrInvalidInput)
	require.ErrorContains(t, err, "no rows to parse")

	assert.Empty(t, actual)
}

func TestParseTable_ReturnsErrorOnNoCells(t *testing.T) {
	t.Parallel()
	// Arrange
	type TestType struct{}

	table := &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{}},
		},
	}

	// Act
	actual, err := ParseTable[TestType](table)

	// Assert
	require.ErrorIs(t, err, ErrInvalidInput)
	require.ErrorContains(t, err, "no cells to parse")

	assert.Empty(t, actual)
}

func TestParseTable_ReturnsErrorOnDecodeFailure(t *testing.T) {
	t.Parallel()
	// Arrange
	type TestType struct {
		Number int
	}

	table := &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{{Value: "number"}}},
			{Cells: []*messages.PickleTableCell{{Value: "123"}}},
			{Cells: []*messages.PickleTableCell{{Value: "456"}}},
		},
	}

	// This will throw an error on Number
	customConfig := &mapstructure.DecoderConfig{WeaklyTypedInput: false}

	// Act
	actual, err := ParseTable[TestType](table, WithDecodeConfig(customConfig))

	// Assert
	require.ErrorContains(t, err, "failed to decode row 0")
	require.ErrorContains(t, err, "failed to decode row 1")
	assert.Empty(t, actual)
}
