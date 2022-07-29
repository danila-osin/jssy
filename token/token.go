package token

type Type struct {
	Name  string
	Regex string
}

var TypesList = []Type{
	{"varDeclaration", "var [a-zA-Z]+"},
	{"printFn", "print:"},
	{"number", "-?[0-9]+([.][0-9]+)?"},
	{"variable", "[a-zA-Z]+"},
	{"semicolon", ";"},
	{"assign", "="},
	{"plus", "\\+"},
	{"minus", "-"},
	{"leftPar", "\\("},
	{"rightPar", "\\)"},
	{"space", "[ \\n\\t\\r]"},
}

func GetTypes(names []string) []Type {
	var tokenTypes []Type

	for _, name := range names {
		for i := range TypesList {
			if TypesList[i].Name == name {
				tokenTypes = append(tokenTypes, TypesList[i])
			}
		}
	}

	return tokenTypes
}

type Token struct {
	Type     Type
	Text     string
	Position int
}

func NewToken(tokenType Type, text string, pos int) *Token {
	return &Token{Type: tokenType, Text: text, Position: pos}
}
