package gowoof

import (
	"testing"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v34"
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

func TestParseTable_ReturnsExpectedDataVertically(t *testing.T) {
	t.Parallel()
	// Arrange
	type TestType struct {
		String  string
		Number  int
		Boolean bool
	}

	table := &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{{Value: "string"}, {Value: "abc"}, {Value: "def"}}},
			{Cells: []*messages.PickleTableCell{{Value: "number"}, {Value: "123"}, {Value: "456"}}},
			{Cells: []*messages.PickleTableCell{{Value: "boolean"}, {Value: "true"}, {Value: "false"}}},
		},
	}

	// Act
	actual, err := ParseTable[TestType](table, Vertical())

	// Assert
	require.NoError(t, err)

	expected := []TestType{
		{String: "abc", Number: 123, Boolean: true},
		{String: "def", Number: 456, Boolean: false},
	}

	assert.Equal(t, expected, actual)
}

func TestParseTable_ReturnsExpectedDataWithNullValues(t *testing.T) {
	t.Parallel()
	// Arrange
	type TestType struct {
		String  *string
		Number  *int
		Boolean *bool
	}

	table := &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{{Value: "string"}, {Value: "number"}, {Value: "boolean"}}},
			{Cells: []*messages.PickleTableCell{{Value: "NULL"}, {Value: "123"}, {Value: "true"}}},
			{Cells: []*messages.PickleTableCell{{Value: "def"}, {Value: "NULL"}, {Value: "false"}}},
			{Cells: []*messages.PickleTableCell{{Value: "ghi"}, {Value: "456"}, {Value: "NULL"}}},
		},
	}

	// Act
	actual, err := ParseTable[TestType](table)

	// Assert
	require.NoError(t, err)

	expected := []TestType{
		{String: nil, Number: new(123), Boolean: new(true)},
		{String: new("def"), Number: nil, Boolean: new(false)},
		{String: new("ghi"), Number: new(456), Boolean: nil},
	}

	assert.Equal(t, expected, actual)
}

func TestParseTable_ReturnsExpectedDataWithPointer(t *testing.T) {
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
	actual, err := ParseTable[*TestType](table)

	// Assert
	require.NoError(t, err)

	expected := []*TestType{
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
	require.ErrorContains(t, err, `generic type "string" is not a struct`)
	assert.Empty(t, actual)
}

func TestParseTable_ReturnsErrorIfGenericIsNotAPointerToAStruct(t *testing.T) {
	t.Parallel()
	// Act
	actual, err := ParseTable[*string](nil)

	// Assert
	require.ErrorIs(t, err, ErrInvalidInput)
	require.ErrorContains(t, err, `generic pointer type "*string" is not a pointer to a struct`)
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
	require.ErrorContains(t, err, "failed to decode")
	assert.Empty(t, actual)
}

func TestTableString_OutputsExpectedTable(t *testing.T) {
	t.Parallel()
	// Arrange
	table := &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{{Value: "string"}, {Value: "number"}, {Value: "boolean"}}},
			{Cells: []*messages.PickleTableCell{{Value: "abc"}, {Value: "123"}, {Value: "true"}}},
			{Cells: []*messages.PickleTableCell{{Value: "def"}, {Value: "123456789"}, {Value: "false"}}},
		},
	}

	// Act
	actual := TableString(table)

	// Assert
	expected := `
| string | number    | boolean |
| abc    | 123       | true    |
| def    | 123456789 | false   |`[1:]

	assert.Equal(t, expected, actual)
}

func TestRowString_OutputsExpectedRow(t *testing.T) {
	t.Parallel()
	// Arrange
	row := &messages.PickleTableRow{
		Cells: []*messages.PickleTableCell{{}, {Value: "def"}, {}, {Value: "123456789"}, {Value: "false"}, {}},
	}

	// Act
	actual := RowString(row)

	// Assert
	expected := `| | def | | 123456789 | false | |`

	assert.Equal(t, expected, actual)
}
