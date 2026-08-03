package gowoof

import (
	"fmt"

	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v34"
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

func verticalGodogTable() *godog.Table {
	return &godog.Table{
		Rows: []*messages.PickleTableRow{
			{Cells: []*messages.PickleTableCell{{Value: "name"}, {Value: "Rex"}, {Value: "Lando"}, {Value: "Bob"}}},
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

func ExampleParseTable_vertically() {
	type Dog struct {
		Name string
	}

	table := verticalGodogTable()

	dogs, err := ParseTable[Dog](table, Vertical())
	if err != nil {
		panic(err)
	}

	fmt.Println(dogs)

	// Output:
	// [{Rex} {Lando} {Bob}]
}
