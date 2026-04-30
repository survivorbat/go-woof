# 🐶 GoWoof

Go Woof has helper functions for godog tests.

## ⬇️ Installation

`go get github.com/survivorbat/go-woof`

## 📋 Usage

### `ParseTable`

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

With the `gowoof.Vertical()` option the function can be instructed to read the table vertically.

## 🔭 Plans

Not much yet.
