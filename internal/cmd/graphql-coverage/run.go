package graphql_coverage

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

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
		re := regexp.MustCompile(`^\s*(#|}|type).*$`)
		var line int
		for scanner.Scan() {
			line++
			coverageCount, _ := coverageMap[line]
			text := scanner.Text()
			coverFile.CoverFile = append(coverFile.CoverFile, coverage.CoverLine{
				Line:  line,
				Count: coverageCount,
				Text:  text,
				Ignored: re.Match([]byte(text)),
			})
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		if err := schemaFile.Close(); err != nil {
			return err
		}

		if *c.HtmlOutput {
			writer, err := os.OpenFile(*c.Output, os.O_RDWR|os.O_CREATE, 0660)
			if err != nil {
				return err
			}
			if err := output.Output(writer, coverFile); err != nil {
				return err
			}
		}

		coverageCount := 0
		for _, coverLine := range coverFile.CoverFile {
			if !coverLine.Ignored && coverLine.Count > 0 {
				coverageCount++
			}
		}

		fmt.Printf("coverage: %.2f%% \n", float64(coverageCount)/float64(len(coverFile.CoverFile))*100)

		return nil
	}
}
