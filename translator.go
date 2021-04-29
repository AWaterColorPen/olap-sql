package olapsql

type Translator interface {
	Translate(interface{}) (interface{}, error)
}
