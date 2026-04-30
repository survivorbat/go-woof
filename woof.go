package gowoof

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/go-viper/mapstructure/v2"
	"github.com/samber/lo"
)

// ErrInvalidInput is returned if an invalid input is encountered
var ErrInvalidInput = errors.New("invalid input")

// ParseTable attempts to parse given table to a []T. T is expected to be a struct. mapstructure is used
// underneath to map the table fields to the output slice.
func ParseTable[T any](table *godog.Table, opts ...Option) ([]T, error) {
	cfg := &Config{
		DecodeConfig: &mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
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
		mapList = parseVertically(table.Rows)
	} else {
		mapList = parseHorizontally(table.Rows)
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

func parseHorizontally(rows []*messages.PickleTableRow) []map[string]string {
	headers := lo.Map(rows[0].Cells, func(cell *messages.PickleTableCell, _ int) string {
		return cell.Value
	})

	mapList := lo.Map(rows[1:], func(row *messages.PickleTableRow, _ int) map[string]string {
		// Not using SliceToMap because the callback lacks the index parameter
		entries := lo.Map(row.Cells, func(item *messages.PickleTableCell, cellIndex int) lo.Entry[string, string] {
			return lo.Entry[string, string]{Key: headers[cellIndex], Value: item.Value}
		})

		return lo.FromEntries(entries)
	})

	return mapList
}

func parseVertically(rows []*messages.PickleTableRow) []map[string]string {
	headers := lo.Map(rows, func(row *messages.PickleTableRow, _ int) string {
		return row.Cells[0].Value
	})

	result := make([]map[string]string, len(rows[0].Cells[1:]))

	for rowIndex, row := range rows {
		for cellIndex, cell := range row.Cells[1:] {
			if result[cellIndex] == nil {
				result[cellIndex] = make(map[string]string, len(rows))
			}

			result[cellIndex][headers[rowIndex]] = cell.Value
		}
	}

	return result
}
