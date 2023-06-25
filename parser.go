package json2json

import (
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"strings"
)

// Parser parses a string and returns the result
type Parser struct {
	input     map[string]interface{}
	funcStack []Func
}

// NewParser creates a new parser
func NewParser(input map[string]interface{}) *Parser {
	return &Parser{
		input: input,
	}
}

// Parse parses a string and returns the result
func (p *Parser) Parse(str string) (any, error) {
	str = p.removeWhitespace(str)
	splitBracketStr := strings.SplitN(str, string(LeftBracket), 2)
	fnStr := Func(splitBracketStr[0])
	args := make([]any, 0)
	_, validFunc := fnFunc[fnStr]
	if validFunc {
		p.funcStack = append(p.funcStack, fnStr)
		defer func() {
			p.funcStack = p.funcStack[:len(p.funcStack)-1]
		}()
		argStr := strings.TrimSuffix(strings.TrimPrefix(str, fmt.Sprintf("%s%s", fnStr, string(LeftBracket))), string(RightBracket))
		argStrSplit := p.splitArgs(argStr)
		for _, a := range argStrSplit {
			arg, err := p.Parse(a)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
	} else {
		if op, ok := containsOp(str); ok {
			idx := strings.LastIndex(str, string(op))
			x, err := p.Parse(str[:idx])
			if err != nil {
				return nil, err
			}
			y, err := p.Parse(str[idx+1:])
			if err != nil {
				return nil, err
			}
			return opFunc[op](x, y), nil
		}
		strInt, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			return strInt, nil
		}
		strFloat, err := strconv.ParseFloat(str, 64)
		if err == nil {
			return strFloat, nil
		}
		if strings.HasPrefix(str, string(Apostrophe)) && strings.HasSuffix(str, string(Apostrophe)) {
			return strings.TrimSuffix(strings.TrimPrefix(str, string(Apostrophe)), string(Apostrophe)), nil
		}
		if strings.HasPrefix(str, string(LeftSquareBracket)) && strings.HasSuffix(str, string(RightSquareBracket)) {
			return p.parseInputByKey(str)
		}
		if constant, ok := constMap[Const(strings.ToUpper(str))]; ok {
			return constant, nil
		}
		return nil, fmt.Errorf("unknown string %s", str)
	}
	res, err := fnFunc[fnStr](args)
	if err != nil {
		return nil, fmt.Errorf("func %s: %w", fnStr, err)
	}
	return res, nil
}

// removeWhitespace remove all spaces
// except if inside two apostrophes which defines a hardcoded string
func (p *Parser) removeWhitespace(str string) string {
	parts := strings.Split(str, string(Apostrophe))
	for i := 0; i < len(parts); i += 2 {
		parts[i] = strings.TrimSpace(parts[i])
	}
	trimmed := strings.Join(parts, string(Apostrophe))
	return trimmed
}

// splitArgs split args string by comma
// except for inside other operators:
// brackets, apostrophes, square brackets
func (p *Parser) splitArgs(str string) []string {
	var result []string
	var buffer strings.Builder
	var openBrackets int
	var openApostrophes bool
	for _, char := range str {
		switch ParserChar(char) {
		case Apostrophe:
			openApostrophes = !openApostrophes
			buffer.WriteRune(char)
		case LeftBracket, LeftSquareBracket:
			openBrackets++
			buffer.WriteRune(char)
		case RightBracket, RightSquareBracket:
			if openBrackets > 0 {
				openBrackets--
				buffer.WriteRune(char)
			}
		case Comma:
			if openBrackets == 0 && !openApostrophes {
				if buffer.Len() > 0 {
					result = append(result, buffer.String())
					buffer.Reset()
				}
			} else {
				buffer.WriteRune(char)
			}
		default:
			buffer.WriteRune(char)
		}
	}
	if buffer.Len() > 0 {
		result = append(result, buffer.String())
	}
	return result
}

// parseInputByKey parse input by key
// e.g. [key1.key2.key3]
func (p *Parser) parseInputByKey(str string) (any, error) {
	key := strings.TrimSuffix(strings.TrimPrefix(str, string(LeftSquareBracket)), string(RightSquareBracket))
	keyParts := strings.Split(key, string(Dot))
	inputTemp := make(map[string]any, len(p.input))
	for k, v := range p.input {
		inputTemp[k] = v
	}
	for i, keyPart := range keyParts {
		if val, ok := inputTemp[keyPart]; ok {
			switch val.(type) {
			case map[string]any:
				if i == len(keyParts)-1 {
					return val, nil
				}
				inputTemp = val.(map[string]any)
			case []any:
				if i == len(keyParts)-1 {
					return val, nil
				}
				inputTemp = make(map[string]any, len(val.([]any)))
				for i, v := range val.([]any) {
					inputTemp[cast.ToString(i)] = v
				}
			case []map[string]any:
				if i == len(keyParts)-1 {
					return val, nil
				}
				inputTemp = make(map[string]any, len(val.([]map[string]any)))
				for i, v := range val.([]map[string]any) {
					inputTemp[cast.ToString(i)] = v
				}
			default:
				return val, nil
			}
		}
	}
	return nil, nil
}
