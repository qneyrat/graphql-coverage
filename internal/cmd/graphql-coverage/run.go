package graphql_coverage

import (
	"bufio"
	"fmt"
	"os"

	"github.com/qneyrat/graphql-coverage/internal/coverage"
	"github.com/qneyrat/graphql-coverage/internal/output"
	"github.com/spf13/cobra"
)

func DecorateRunFunc(c *Context) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		calc := coverage.NewCalculator(*c.IsJson)
		err := calc.LoadSchema(string(c.SchemaContent))
		if err != nil {
			return err
		}

		coverageMap, err := calc.Calculate(c.InputFiles)
		if err != nil {
			return err
		}

		schemaFile, err := os.Open(*c.SchemaFilename)
		if err != nil {
			return err
		}

		coverFile := coverage.WrappedCoverFile{[]coverage.CoverLine{}}
		scanner := bufio.NewScanner(schemaFile)
		var count int
		for scanner.Scan() {
			count++
			coverageCount, ok := coverageMap[count]
			if !ok {
				coverFile.CoverFile = append(coverFile.CoverFile, coverage.CoverLine{
					Line:  count,
					Text:  scanner.Text(),
				})
				continue
			}

			coverFile.CoverFile = append(coverFile.CoverFile, coverage.CoverLine{
				Line:  count,
				Count: coverageCount,
				Text:  scanner.Text(),
			})
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		if err := schemaFile.Close(); err != nil {
			return err
		}

		if *c.HtmlOutput {
			if err := output.Output(*c.Output, coverFile); err != nil {
				return err
			}
		}

		for _, coverLine := range coverFile.CoverFile {
			fmt.Println(coverLine.Line, coverLine.Text, coverLine.Count)
		}

		return nil
	}
}
