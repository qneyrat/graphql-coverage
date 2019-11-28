package coverage_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qneyrat/graphql-coverage/internal/coverage"
)

var schema string
var queries = map[string][]byte{}
var outs = map[string][]byte{}

func TestMain(m *testing.M) {
	file, err := os.Open("testdata/schema.graphql")
	if err != nil {
		fmt.Println("error on Open")
		os.Exit(1)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error on ReadAll")
		os.Exit(1)
	}

	_ = file.Close()
	schema = string(content)
	err = filepath.Walk("testdata/queries", func(path string, info os.FileInfo, err error) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(path, filepath.Ext(path))
		switch filepath.Ext(path) {
		case ".graphql":
			content, _ := ioutil.ReadAll(file)
			_ = file.Close()
			queries[name] = content
			break

		case ".out":
			content, _ := ioutil.ReadAll(file)
			_ = file.Close()
			outs[name] = content
		}

		return nil
	}); if err != nil {
		fmt.Println("error on Walk")
		os.Exit(1)
	}

	result := m.Run()
	os.Exit(result)
}

func TestCalculator_Calculate(t *testing.T) {
	calc := coverage.NewCalculator(false)
	err := calc.LoadSchema(schema)
	if err != nil {
		t.Error("LoadSchema should return nil, actual:", err.Error())
	}

	for identifier, query := range queries{
		buf := bytes.NewBuffer(query)
		coverageMap, err := calc.Calculate([]io.ReadCloser{ioutil.NopCloser(buf)})
		if err != nil {
			t.Error("Calculate should return nil, actual:", err.Error())
		}

		expectedCoverage := outs[identifier]
		actualCoverage, err := json.Marshal(coverageMap)
		if err != nil {
			t.Error("json.Marshal should return nil, actual:", err.Error())
		}

		expectedCoverageStr := string(expectedCoverage)
		actualCoverageStr := string(actualCoverage)

		if expectedCoverageStr != actualCoverageStr {
			t.Error("Expected:", expectedCoverageStr, "actual:", actualCoverageStr)
		}
	}
}
