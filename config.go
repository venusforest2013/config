package config

const (
	TomlType = Type("Toml")
	JsonType = Type("Json")

	Delimiter = "."
)

type Type string

type Node interface {
	Type() Type
	Value() interface{}
	Access(string) Node
	AccessArray(string) []Node
	AccessMap(string) map[string]Node
	String() string
}

//type Value interface {
//	Node
//
//	Has(string) bool
//	HasArray(string) bool
//	HasMap(string) bool
//
//	Bool(string, bool) bool
//	Str(string, string) string
//	Int64(string, int64) int64
//	Float64(string, float64) float64
//	Duration(string, time.Duration) time.Duration
//
//	BoolArray(string) []bool
//	StrArray(string) []string
//	Int64Array(string) []int64
//	Float64Array(string) []float64
//
//	ValueArray() []*Value
//	ValueMap() map[string]*Value
//}
