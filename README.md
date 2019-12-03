# graphql-coverage

```
Calculate coverage on GraphQL schema

Usage:
  graphql-coverage [graphql query file or dir] [flags]

Flags:
  -f, --filter string   filter to search graphql queries when arg is dir (default ".graphql")
  -h, --help            help for graphql-coverage
      --html            html output
      --json            json input graphql query
  -m, --max-depth int   max depth to search graphql queries when arg is dir (default 1)
  -o, --output string   output file (default "coverage.out")
  -s, --schema string   graphql schema (default "schema.graphql")
```                   

## Example usage

### Basic:
```
go run cmd/graphql-coverage/main.go --schema examples/schema.graphql examples/queries
```

### HTML output
```
go run cmd/graphql-coverage/main.go --schema examples/schema.graphql --html --output coverage.html examples/queries
```

### Input queries
Use `-f` to search graphql queries in path (by default `.graphql`)

Use `--json` if your queries are graphql queries in json like 
```
{
    "query":"query YourQuery) {}",
    "variables":{"yourvar":"value"},
    "operationName":"YourQuery"
}
```
                                                                                        