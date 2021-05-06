package types

type Statement interface {
	Statement() (string, error)
}
