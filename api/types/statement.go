package types

type Statement interface {
	Expression() (string, error) // Any expression
	Alias() (string, error)      // Name for expr. Aliases should comply with the identifiers syntax
	Statement() (string, error)  // expr AS alias
}
