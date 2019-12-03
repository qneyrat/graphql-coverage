package graphql_coverage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func DecoratePreRunFunc(c *Context) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if c.IsDir {
				err := filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					// TODO Use maxDepth and recursive func
					return nil
				}

				// TODO use Regex
				if filepath.Ext(path) != *c.Filter {
					return nil
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}

				c.InputFiles = append(c.InputFiles, file)

				return nil
			}); if err != nil {
				return err
			}
		} else {
			// TODO Merge code with isDir part
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}

			c.InputFiles = append(c.InputFiles, file)

		}

		schemaFile, err := os.Open(*c.SchemaFilename)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		c.SchemaContent, err = ioutil.ReadAll(schemaFile)
		if err != nil {
			return err
		}

		if err = schemaFile.Close();  err != nil {
			return err
		}

		return nil
	}
}
