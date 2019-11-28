package graphql_coverage

import (
	"io"
)

type Context struct {
	SchemaFilename *string
	Output *string
	HtmlOutput *bool
	MaxDepth *int
	Filter *string
	IsJson *bool
	Debug *bool

	IsDir bool
	InputFiles []io.ReadCloser
	SchemaContent []byte

}

type Option func(c *Context)
