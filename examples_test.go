package gowoof

import (
	"fmt"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
)

func godogTable() *godog.Table {
	return &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{{Value: "name"}}},
			{Cells: []*messages.PickleTableCell{{Value: "Rex"}}},
			{Cells: []*messages.PickleTableCell{{Value: "Lando"}}},
			{Cells: []*messages.PickleTableCell{{Value: "Bob"}}},
		},
	}
}

func ExampleParseTable() {
	type Dog struct {
		Name string
	}

	table := godogTable()

	dogs, err := ParseTable[Dog](table)
	if err != nil {
		panic(err)
	}

	fmt.Println(dogs)

	// Output:
	// [{Rex} {Lando} {Bob}]
}
