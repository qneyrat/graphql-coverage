package coverage

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
)

type GraphqlQuery struct {
	Query     string `json:"query"`
	Variables map[string]interface{} `json:"variables"`
	OperationName string `json:"operationName"`
}

func getLine(position *ast.Position) int {
	if position == nil {
		return -1
	}

	return position.Line
}

func coverOnSelectionSet(coverage map[int]int, selectionSet ast.SelectionSet) {
	for _, selection := range selectionSet {
		if field, ok := selection.(*ast.Field); ok {
			coverage[getLine(field.Definition.Position)]++
			for _, arg := range field.Arguments {
				if arg.Value.Definition != nil && arg.Value.Definition.EnumValues.ForName(arg.Value.Raw) != nil {
					coverage[getLine(arg.Value.Definition.EnumValues.ForName(arg.Value.Raw).Position)]++
					continue
				}

				coverage[getLine(arg.Position)]++
			}

			if len(field.SelectionSet) > 0 {
				coverOnSelectionSet(coverage, field.SelectionSet)
			}
		}
	}
}

type WrappedCoverFile struct {
	CoverFile []CoverLine
}

type CoverLine struct {
	Line int
	Count int
	Text string
	Ignored bool
}

type Calculator struct {
	Schema *ast.Schema

	JsonQueries bool
}

func NewCalculator(jsonQueries bool) *Calculator {
	return &Calculator{
		JsonQueries: jsonQueries,
	}
}

func (c *Calculator) LoadSchema(schema string) error {
	parsedSchema, gqlErr := gqlparser.LoadSchema(&ast.Source{
		Name: "coverage",
		Input: schema,
	})
	if gqlErr != nil {
		return gqlErr
	}

	c.Schema = parsedSchema

	return nil
}

func (c *Calculator) Calculate(queries []io.ReadCloser) (map[int]int, error) {
	coverageMap := map[int]int{}
	for _, inputQuery := range queries {
		inputFileContent, err := ioutil.ReadAll(inputQuery)
		_ = inputQuery.Close()
		if err != nil {
			return nil, err
		}

		var query string
		if c.JsonQueries {
			var completeQuery GraphqlQuery
			err = json.Unmarshal(inputFileContent, &completeQuery)
			if err != nil {
				return nil, err
			}
			query = completeQuery.Query
		} else {
			query = string(inputFileContent)
		}

		parsedQuery := gqlparser.MustLoadQuery(c.Schema, query)
		for _, operation := range parsedQuery.Operations {
			coverOnSelectionSet(coverageMap, operation.SelectionSet)
		}
	}

	return coverageMap, nil
}
