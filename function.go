package json2json

import (
	"fmt"
	"github.com/spf13/cast"
	"golang.org/x/exp/utf8string"
	"math"
	"reflect"
	"unicode/utf8"
)

// Func is a function
type Func string

const (
	String   Func = "STRING"
	Int      Func = "INT"
	Float    Func = "FLOAT"
	Bool     Func = "BOOL"
	Object   Func = "OBJECT"
	Array    Func = "ARRAY"
	Var      Func = "VAR"
	Set      Func = "SET"
	Len      Func = "LEN"
	SliceStr Func = "SLICE_STR"
	If       Func = "IF"
	Switch   Func = "SWITCH"
	And      Func = "AND"
	Or       Func = "OR"
	Gte      Func = "GTE"
	Gt       Func = "GT"
	Lte      Func = "LTE"
	Lt       Func = "LT"
)

// funcMap is a map that contains all functions
var fnFunc = map[Func]func([]any) (any, error){
	String:   stringFunc,
	Int:      intFunc,
	Float:    floatFunc,
	Bool:     boolFunc,
	Object:   objectFunc,
	Array:    arrayFunc,
	Var:      varFunc,
	Set:      setFunc,
	Len:      lenFunc,
	SliceStr: sliceStrFunc,
	If:       ifFunc,
	Switch:   switchFunc,
	And:      andFunc,
	Or:       orFunc,
	Gte:      gteFunc,
	Gt:       gtFunc,
	Lte:      lteFunc,
	Lt:       ltFunc,
}

// stringFunc is the string function
// STRING(expr)
// convert expr to string
func stringFunc(args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	return cast.ToStringE(args[0])
}

// intFunc is the int function
// INT(expr)
// convert expr to int
func intFunc(args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	return cast.ToInt64E(args[0])
}

// floatFunc is the float function
// FLOAT(expr, precision)
// convert expr to float
// precision is the number of digits after the decimal point
// default precision is 2
func floatFunc(args []any) (any, error) {
	switch len(args) {
	case 1, 2:
		floatNum, err := cast.ToFloat64E(args[0])
		if err != nil {
			return nil, err
		}
		precision := 2.0
		if len(args) == 2 {
			uintPrecision := uint(2)
			uintPrecision, err = cast.ToUintE(args[1])
			if err != nil {
				return nil, err
			}
			precision, _ = cast.ToFloat64E(uintPrecision)
		}
		precisionPow := math.Pow(10, precision)
		return math.Round(floatNum*precisionPow) / precisionPow, nil
	default:
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
}

// boolFunc is the bool function
// BOOL(expr)
// convert expr to bool
func boolFunc(args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	return cast.ToBoolE(args[0])
}

// objectFunc is the object function
// OBJECT(expr, default)
// if expr is a valid, return true, else return default
func objectFunc(args []any) (any, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	valid := cast.ToBool(args[0])
	if !valid {
		if len(args) == 2 {
			return args[1], nil
		}
	}
	return valid, nil
}

// arrayFunc is the array function
// ARRAY(expr, default)
// if expr is a valid array, return true, else return default
func arrayFunc(args []any) (any, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	_, ok := args[0].([]any)
	if !ok {
		_, ok = args[0].([]map[string]any)
		if !ok {
			if len(args) == 2 {
				return args[1], nil
			}
		}
	}
	return true, nil
}

// setFunc is the set function
// SET(expr)
// return expr
func setFunc(args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	return args[0], nil
}

// varFunc is the var function
// VAR(expr, default)
// if expr is nil, return default, else return expr
func varFunc(args []any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	if args[0] == nil {
		return args[1], nil
	}
	return args[0], nil
}

// lenFunc is the len function
// LEN(expr)
// return the length of expr
// if expr is a string, return the number of characters
// if expr is an array, return the number of elements
// if expr is an object, return the number of keys
// else return error
func lenFunc(args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	switch arg := args[0].(type) {
	case string:
		return utf8.RuneCountInString(arg), nil
	case []any:
		return len(arg), nil
	case map[string]any:
		return len(arg), nil
	default:
		return nil, fmt.Errorf("invalid type: %T", arg)
	}
}

// sliceStrFunc is the sliceStr function
// SLICE_STR(str, start, end)
// return the substring of str from start to end
func sliceStrFunc(args []any) (anyStr any, err error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	str, _ := cast.ToStringE(args[0])
	start, err := cast.ToIntE(args[1])
	if err != nil {
		return
	}
	end, err := cast.ToIntE(args[2])
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid arguments: %v", r)
		}
	}()
	anyStr = utf8string.NewString(str).Slice(start, end)
	return
}

// ifFunc is the if function
// IF(expr, x, y)
// if expr is true, return x, else return y
func ifFunc(args []any) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	expr, err := cast.ToBoolE(args[0])
	if err != nil {
		return nil, err
	}
	if expr {
		return args[1], nil
	}
	return args[2], nil
}

// switchFunc is the switch function
// SWITCH(expr, x1, y1, x2, y2, ..., xn, yn, default)
// if expr == x1, return y1
// if expr == x2, return y2
// ...
// if expr == xn, return yn
// else return default
func switchFunc(args []any) (any, error) {
	if (len(args) % 2) != 0 {
		return nil, fmt.Errorf("invalid odd number of arguments: %d", len(args))
	}
	if len(args) == 0 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	if len(args) == 2 {
		return args[1], nil
	}
	expr := args[0]
	for i := 1; i < len(args); i += 2 {
		if reflect.TypeOf(expr) == reflect.TypeOf(args[i]) {
			if reflect.DeepEqual(expr, args[i]) {
				return args[i+1], nil
			}
		}
	}
	return args[len(args)-1], nil
}

// andFunc is the and function
// AND(expr1, expr2, ..., exprn)
// if all expr are true, return true, else return false
func andFunc(args []any) (any, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	for _, arg := range args {
		if expr, err := cast.ToBoolE(arg); err != nil {
			return false, err
		} else if !expr {
			return false, nil
		}
	}
	return true, nil
}

// orFunc is the or function
// OR(expr1, expr2, ..., exprn)
// if any expr is true, return true, else return false
func orFunc(args []any) (any, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	for _, arg := range args {
		if expr, err := cast.ToBoolE(arg); err != nil {
			return false, err
		} else if expr {
			return true, nil
		}
	}
	return false, nil
}

// gteFunc is the greater than or equal function
// GTE(expr1, expr2)
// if expr1 >= expr2, return true, else return false
func gteFunc(args []any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	first, err := cast.ToFloat64E(args[0])
	if err != nil {
		return nil, err
	}
	second, err := cast.ToFloat64E(args[1])
	if err != nil {
		return nil, err
	}
	return first >= second, nil
}

// gtFunc is the greater than function
// GT(expr1, expr2)
// if expr1 > expr2, return true, else return false
func gtFunc(args []any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	first, err := cast.ToFloat64E(args[0])
	if err != nil {
		return nil, err
	}
	second, err := cast.ToFloat64E(args[1])
	if err != nil {
		return nil, err
	}
	return first > second, nil
}

// lteFunc is the less than or equal function
// LTE(expr1, expr2)
// if expr1 <= expr2, return true, else return false
func lteFunc(args []any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	first, err := cast.ToFloat64E(args[0])
	if err != nil {
		return nil, err
	}
	second, err := cast.ToFloat64E(args[1])
	if err != nil {
		return nil, err
	}
	return first <= second, nil
}

// ltFunc is the less than function
// LT(expr1, expr2)
// if expr1 < expr2, return true, else return false
func ltFunc(args []any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: %d", len(args))
	}
	first, err := cast.ToFloat64E(args[0])
	if err != nil {
		return nil, err
	}
	second, err := cast.ToFloat64E(args[1])
	if err != nil {
		return nil, err
	}
	return first < second, nil
}
