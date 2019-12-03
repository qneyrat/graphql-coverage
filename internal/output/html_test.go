package output_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/qneyrat/graphql-coverage/internal/coverage"
	"github.com/qneyrat/graphql-coverage/internal/output"
)

func TestOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	data := coverage.WrappedCoverFile{
		CoverFile: []coverage.CoverLine{
			{
				Line:  1,
				Count: 0,
				Text:  "# The mutation type, represents all updates we can make to our data",
			},
			{
				Line:  2,
				Count: 1,
				Text:  "type Mutation {",
			},
			{
				Line:  3,
				Count: 0,
				Text:  "}",
			},
		},
	}
	err := output.Output(buf, data)
	if err != nil {
		t.Error("Output should return nil, actual:", err.Error())
	}

	fmt.Println(buf.String())
}

func TestRegex(t *testing.T) {
	data := map[string]bool{}
	data["# aaaa"] = true
	data["type foo {"] = true
	data["	aaaa {"] = false
	data["	} aaa"] = true
	data["program: Program"] = false
	data["	# Liste d'articles"] = true
	data["	# Décoration associée aux liens freemium"] = true

	re := regexp.MustCompile(`^\s*(#|}|type).*$`)
	for line, match := range data {
		actualMatch := re.Match([]byte(line))
		if actualMatch != match {
			t.Error("Match not equal, actual:", actualMatch, "expected:", match)
		}
	}
}
