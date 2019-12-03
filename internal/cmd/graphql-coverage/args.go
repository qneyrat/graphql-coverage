package graphql_coverage

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func DecorateArgsFunc(c *Context) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg, received %d", len(args))
		}

		fi, err := os.Stat(args[0])
		if err != nil {
			return err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			c.IsDir = true
		case mode.IsRegular():
			c.IsDir = false
		default:
			return fmt.Errorf("file or dir %s doesnt exist", args[0])
		}

		return nil
	}
}
