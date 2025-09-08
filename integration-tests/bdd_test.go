package integrationtests

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/survivorbat/go-gowoof"
)

type Dog struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Type       string `json:"type"`
	Vaccinated bool   `json:"Vaccinated"`
}

type scenario struct {
	Actual []Dog
}

func (s *scenario) iHaveAStructTypeThatLooksLikeTheFollowingStructure(*godog.DocString) error {
	// See above
	return nil
}

func (s *scenario) iUseTheFromTableFunctionWithTheTable(ctx context.Context, input *godog.Table) error {
	t := godog.T(ctx)

	actual, err := gowoof.ParseTable[Dog](input)
	require.NoError(t, err)

	s.Actual = actual

	return nil
}

func (s *scenario) iExpectASliceThatResemblesTheFollowingJSON(ctx context.Context, input *godog.DocString) error {
	t := godog.T(ctx)

	var expected []Dog

	err := json.Unmarshal([]byte(input.Content), &expected)
	require.NoError(t, err)

	assert.Equal(t, expected, s.Actual)

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	scenario := new(scenario)

	ctx.Step(`^I expect a slice that resembles the following JSON:$`, scenario.iExpectASliceThatResemblesTheFollowingJSON)
	ctx.Step(`^I have a struct type that looks like the following structure:$`, scenario.iHaveAStructTypeThatLooksLikeTheFollowingStructure)
	ctx.Step(`^I use the FromTable function with the table:$`, scenario.iUseTheFromTableFunctionWithTheTable)
}

func TestFeatures(t *testing.T) {
	t.Parallel()

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
