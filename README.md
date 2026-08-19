# 🐶 GoWoof

Go Woof has helper functions for godog tests.

## ⬇️ Installation

`go get github.com/survivorbat/go-woof`

## 📋 Usage

### ParseTable

Parse a table to a struct using `mapstructure`.

```go
package tests

import (
  "github.com/cucumber/godog"
  "github.com/survivorbat/go-woof"
  "github.com/stretchr/testify/require"
)

type Dog struct {
  Name string
}

type scenario struct {
  Dogs []Dog
}

func (s *scenario) theFollowingDogsAreInTheSystem(ctx context.Context, dogTable *godog.Table) error {
  t := godog.T(ctx)

  dogs, err := gowoof.ParseTable[Dog](dogTable)
  require.NoError(t, err)

  s.Dogs = dogs

  return nil
}
```

#### Options

- `gowoof.Vertical()` will cause the table to be read vertically (first column as labels)
- `gowoof.WithNullValue("...")` will set the string used to indicate that a cell should be skipped, defaults to `NULL`

### TableString

Stringify a godog table into a gherkin-like format.

```go
fmt.Println(gowoof.TableString(godogTable))
```

### RowString

Stringify a godog row into a gherkin-like format.

```go
fmt.Println(gowoof.RowString(godogRow))
```

## 🔭 Plans

Not much yet.
