package json2json

import (
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

// Operator is an operator
type Operator string

const (
	Eq    Operator = "="
	NotEq Operator = "<>"
	Mul   Operator = "*"
	Div   Operator = "/"
	Add   Operator = "+"
	Sub   Operator = "-"
)

// opFunc is a map that contains all operator functions
var opFunc = map[Operator]func(x, y any) any{
	Eq:    eqFunc,
	NotEq: notEqFunc,
	Mul:   mulFunc,
	Div:   divFunc,
	Add:   addFunc,
	Sub:   subFunc,
}

// containsOp checks if a string contains an operator
func containsOp(str string) (Operator, bool) {
	for op := range opFunc {
		if strings.Contains(str, string(op)) {
			return op, true
		}
	}
	return "", false
}

// eqFunc checks if two values are equal
func eqFunc(x, y any) any {
	return reflect.DeepEqual(x, y)
}

// notEqFunc checks if two values are not equal
func notEqFunc(x, y any) any {
	return !reflect.DeepEqual(x, y)
}

// mulFunc multiplies two values
func mulFunc(x, y any) any {
	return cast.ToFloat64(x) * cast.ToFloat64(y)
}

// divFunc divides two values
func divFunc(x, y any) any {
	return cast.ToFloat64(x) / cast.ToFloat64(y)
}

// addFunc adds two values
func addFunc(x, y any) any {
	return cast.ToFloat64(x) + cast.ToFloat64(y)
}

// subFunc subtracts two values
func subFunc(x, y any) any {
	return cast.ToFloat64(x) - cast.ToFloat64(y)
}
