package graphql_coverage

import (
	flag "github.com/spf13/pflag"
)

func WithFlags(flags *flag.FlagSet) *Context {
	return &Context{
		SchemaFilename: flags.StringP(
			"schema",
			"s",
			"schema.graphql",
			"graphql schema",
		),
		Output: flags.StringP(
			"output",
			"o",
			"coverage.out",
			"output file",
		),
		MaxDepth: flags.IntP(
			"max-depth",
			"m",
			1,
			"max depth to search graphql queries when arg is dir",
		),
		Filter: flags.StringP(
			"filter",
			"f",
			".graphql",
			"filter to search graphql queries when arg is dir",
		),
		HtmlOutput: flags.Bool("html", false, "html output"),
		IsJson: flags.Bool("json", false, "json input graphql query"),
		IsDir: false,
		InputFiles: nil,
		SchemaContent: nil,
	}
}
