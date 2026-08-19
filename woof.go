package gowoof

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v34"
	"github.com/go-viper/mapstructure/v2"
)

// ErrInvalidInput is returned if an invalid input is encountered
var ErrInvalidInput = errors.New("invalid input")

// ParseTable attempts to parse given table to a []T. T is expected to be a struct. mapstructure is used
// underneath to map the table fields to the output slice.
func ParseTable[T any](table *godog.Table, opts ...Option) ([]T, error) {
	cfg := &Config{
		NullValue: "NULL",
		DecodeConfig: &mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Squash:           true,
		},
	}

	for index, option := range opts {
		err := option(cfg)
		if err != nil {
			return nil, fmt.Errorf("option %d returned an error: %w", index, err)
		}
	}

	var t T

	tType := reflect.TypeOf(t)

	if tType.Kind() != reflect.Struct && tType.Kind() != reflect.Pointer {
		return nil, fmt.Errorf(`generic type "%T" is not a struct or pointer to a struct: %w`, t, ErrInvalidInput)
	} else if tType.Kind() == reflect.Pointer && tType.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf(`generic pointer type "%T" is not a pointer to a struct: %w`, t, ErrInvalidInput)
	}

	// No rows is _not_ an empty table, we need a header line
	if len(table.Rows) == 0 {
		return nil, fmt.Errorf("no rows to parse: %w", ErrInvalidInput)
	} else if len(table.Rows[0].Cells) == 0 {
		return nil, fmt.Errorf("no cells to parse: %w", ErrInvalidInput)
	}

	var mapList []map[string]string

	if cfg.Vertical {
		mapList = parseVertically(cfg, table.Rows)
	} else {
		mapList = parseHorizontally(cfg, table.Rows)
	}

	var result []T

	// Required for the function to work
	cfg.DecodeConfig.Result = &result

	decoder, err := mapstructure.NewDecoder(cfg.DecodeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate decoder: %w", err)
	}

	// Use mapstructure to save us from the mapping/parsing work
	err = decoder.Decode(mapList)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return result, nil
}

// parseHorizontally takes the first row as a list of headers, and every row after that as entries
func parseHorizontally(cfg *Config, rows []*messages.PickleTableRow) []map[string]string {
	headers := make([]string, len(rows[0].Cells))
	for index, cell := range rows[0].Cells {
		headers[index] = cell.Value
	}

	output := make([]map[string]string, len(rows[1:]))

	for rowIndex, row := range rows[1:] {
		output[rowIndex] = make(map[string]string, len(rows[0].Cells))

		for cellIndex, cell := range row.Cells {
			if cell.Value == cfg.NullValue {
				continue
			}

			output[rowIndex][headers[cellIndex]] = cell.Value
		}
	}

	return output
}

// parseVertically takes the first column as a list of headers, and every column after that as entries
func parseVertically(cfg *Config, rows []*messages.PickleTableRow) []map[string]string {
	headers := make([]string, len(rows[0].Cells))
	for index, row := range rows {
		headers[index] = row.Cells[0].Value
	}

	output := make([]map[string]string, len(rows[0].Cells[1:]))

	for rowIndex, row := range rows {
		for cellIndex, cell := range row.Cells[1:] {
			if cell.Value == cfg.NullValue {
				continue
			}

			if output[cellIndex] == nil {
				output[cellIndex] = make(map[string]string, len(rows))
			}

			output[cellIndex][headers[rowIndex]] = cell.Value
		}
	}

	return output
}

// TableString turns the table into a gherkin-like representation for output purposes
func TableString(table *godog.Table) string {
	columnWidths := make(map[int]int, len(table.Rows[0].Cells))

	for _, row := range table.Rows {
		for index, cell := range row.Cells {
			if columnWidths[index] < len(cell.Value) {
				columnWidths[index] = len(cell.Value)
			}
		}
	}

	var builder strings.Builder
	for rowIndex, row := range table.Rows {
		_, _ = builder.WriteRune('|')

		for cellIndex, cell := range row.Cells {
			_, _ = fmt.Fprintf(&builder, " %-*s |", columnWidths[cellIndex], cell.Value)
		}

		if rowIndex < len(table.Rows)-1 {
			_, _ = builder.WriteRune('\n')
		}
	}

	return builder.String()
}

// RowString turns the row into a gherkin-like representation for output purposes
func RowString(row *messages.PickleTableRow) string {
	values := make([]string, len(row.Cells))
	for index, cell := range row.Cells {
		values[index] = cell.Value
	}

	return "| " + strings.Join(values, " | ") + " |"
}
