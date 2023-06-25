package json2json

// Const is a constant
type Const string

const (
	Nil        Const = "NIL"
	True       Const = "TRUE"
	False      Const = "FALSE"
	NoParam    Const = "NO_PARAM"
	EmptyArray Const = "EMPTY_ARRAY"
)

// NoParamVar is a variable that represents no parameter
var NoParamVar any = nil

// constMap is a map that contains all constants
var constMap = map[Const]any{
	Nil:        nil,
	True:       true,
	False:      false,
	NoParam:    &NoParamVar,
	EmptyArray: []any{},
}
