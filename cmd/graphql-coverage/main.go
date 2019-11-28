package main

import (
	"os"

	graphqlCoverage "github.com/qneyrat/graphql-coverage/internal/cmd/graphql-coverage"

)

func main() {
	err := graphqlCoverage.NewCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
