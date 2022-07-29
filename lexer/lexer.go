package lexer

import (
	"fmt"
	tok "github.com/strange-cat-fe/jssy/token"
	"log"
	"regexp"
)

type Lexer struct {
	Code      string
	Position  int
	TokenList []tok.Token
}

func NewLexer(code string) *Lexer {
	return &Lexer{Code: code}
}

func (l *Lexer) LexicalAnalysis() []tok.Token {
	for l.nextToken() {}

	var noSpacesList []tok.Token

	for _, token := range l.TokenList {
		if token.Text != " " && token.Text != "\n" {
			noSpacesList = append(noSpacesList, token)
		}
	}

	fmt.Println(noSpacesList)

	return noSpacesList
}

func (l *Lexer) nextToken() bool {
	if l.Position >= len(l.Code) {
		return false
	}

	tokenTypes := tok.TypesList

	for _, tokenType := range tokenTypes {
		regex, err := regexp.Compile("^" + tokenType.Regex)
		if err != nil {
			log.Fatalf("Can't compile regexp for token %s", tokenType.Name, err.Error())
		}

		tokenStr := regex.FindString(l.Code[l.Position:])

		if tokenStr != "" {
			token := *tok.NewToken(tokenType, tokenStr, l.Position)
			l.TokenList = append(l.TokenList, token)

			l.Position += len(tokenStr)
			return true
		}
	}

	log.Fatalf("There is an syntax error on position: %d", l.Position)
	return false
}
