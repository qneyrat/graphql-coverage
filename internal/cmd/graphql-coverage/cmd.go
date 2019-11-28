package graphql_coverage

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "graphql-coverage [graphql query file or dir]",
		Short: "Calculate coverage on GraphQL schema",
	}

	cmdCtx := WithFlags(rootCmd.PersistentFlags())
	rootCmd.Args = DecorateArgsFunc(cmdCtx)
	rootCmd.PreRunE = DecoratePreRunFunc(cmdCtx)
	rootCmd.RunE = DecorateRunFunc(cmdCtx)

	return rootCmd
}
