package json2json

// ParserChar is a type for the characters
// that are used in the parser
type ParserChar string

const (
	// LeftBracket is the left bracket
	// for opening the function call arguments
	LeftBracket ParserChar = "("
	// RightBracket is the right bracket
	// for closing the function call arguments
	RightBracket ParserChar = ")"

	// LeftSquareBracket is the left square bracket
	// for opening the input key path
	LeftSquareBracket ParserChar = "["
	// RightSquareBracket is the right square bracket
	// for closing the input key path
	RightSquareBracket ParserChar = "]"

	// Apostrophe is the apostrophe
	// for defining a string
	Apostrophe ParserChar = "'"

	// Dot is the dot
	// for defining a key path inside the square brackets
	Dot ParserChar = "."

	// Comma is the comma
	// for separating the function call arguments
	Comma ParserChar = ","
)
