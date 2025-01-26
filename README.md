# Query Builder (qbr)

`qbr` is a flexible and lightweight multi-purpose query builder for constructing various types of queries. Currently, it supports `SQL` queries, but it is designed to be extensible for other query types in the future. It allows you to build SELECT, INSERT, UPDATE, and DELETE queries with complex conditions, sorting, and pagination.

## Installation

To install `qbr`, use the following command:

```bash
go get -u github.com/tyrenix/goqbr
```

## Quick Start

Here's a simple example of how to use `qbr` to build a SELECT query:

```go
package main

import (
    "fmt"
    "github.com/tyrenix/qbr"
)

func main() {
    qb := qbr.
        New(qbr.QueryTypeSelect).
        Select(
            qbr.NewField("name", qbr.FieldDefault),
        ).
        Where(
            qbr.Eq(
                qbr.NewField("age", qbr.FieldDefault), 
                25,
            ),
        ).
        Limit(100).
        Offset(50)

    query, params, err := qb.ToSql("users", qbr.ToSqlDollar)
    if err != nil {
        panic(err)
    }

    fmt.Println("Query:", query)
    fmt.Println("Params:", params)
}
```

## Documentation

For detailed documentation, including API reference and additional examples, visit the official documentation on [pkg.go.dev](https://pkg.go.dev/github.com/tyrenix/goqbr).

## Roadmap

* **Current Features:**
    * Support for SQL queries (SELECT, INSERT, UPDATE, DELETE).
    * Complex conditions, sorting, and pagination.
* **Future Plans:**
    * Extend support for other query types (e.g., NoSQL, GraphQL).
    * Add more advanced query building features.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/tyrenix/goqbr/blob/master/LICENSE) file for details.