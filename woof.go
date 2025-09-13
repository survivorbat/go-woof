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

	// Get all the fields/headers
	headers := lo.Map(table.Rows[0].Cells, func(item *messages.PickleTableCell, _ int) string {
		return item.Value
	})

	var result []T

	mapList := lo.Map(table.Rows[1:], func(row *messages.PickleTableRow, rowIndex int) map[string]string {
		// Not using SliceToMap because the callback lacks the index parameter
		entries := lo.Map(row.Cells, func(item *messages.PickleTableCell, cellIndex int) lo.Entry[string, string] {
			return lo.Entry[string, string]{Key: headers[cellIndex], Value: item.Value}
		})

		return lo.FromEntries(entries)
	})

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
