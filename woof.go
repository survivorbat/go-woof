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

var ErrInvalidInput = errors.New("invalid input")

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

	if reflect.ValueOf(*new(T)).Kind() != reflect.Struct {
		return nil, fmt.Errorf("generic type is not a struct: %w", ErrInvalidInput)
	}

	if len(table.Rows) == 0 {
		return nil, fmt.Errorf("no rows to parse: %w", ErrInvalidInput)
	} else if len(table.Rows[0].Cells) == 0 {
		return nil, fmt.Errorf("no cells to parse: %w", ErrInvalidInput)
	}

	errs := make([]error, len(table.Rows)-1)

	// Get all the fields/headers
	headers := lo.Map(table.Rows[0].Cells, func(item *messages.PickleTableCell, _ int) string {
		return item.Value
	})

	result := lo.Map(table.Rows[1:], func(row *messages.PickleTableRow, rowIndex int) T {
		var item T

		// Not using SliceToMap because the callback lacks the index parameter
		mapValues := lo.Map(row.Cells, func(item *messages.PickleTableCell, cellIndex int) lo.Entry[string, string] {
			return lo.Entry[string, string]{Key: headers[cellIndex], Value: item.Value}
		})

		// Required for the function to work
		cfg.DecodeConfig.Result = &item

		decoder, err := mapstructure.NewDecoder(cfg.DecodeConfig)
		if err != nil {
			errs[rowIndex] = fmt.Errorf("failed to instantiate decoder: %w", err)
			return item
		}

		// Use mapstructure to save us from the mapping/parsing work
		err = decoder.Decode(lo.FromEntries(mapValues))
		if err != nil {
			errs[rowIndex] = fmt.Errorf("failed to decode row %d: %w", rowIndex, err)
		}

		return item
	})

	err := errors.Join(errs...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
