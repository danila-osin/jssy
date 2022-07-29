package main

import (
	"io/ioutil"
	"log"

	lex "github.com/strange-cat-fe/jssy/lexer"
	par "github.com/strange-cat-fe/jssy/parser"
)

func readFile(path string) string {

	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {
	code := readFile("fuck.jssy")

	lexer := lex.NewLexer(code)

	tokenList := lexer.LexicalAnalysis()

	parser := par.NewParser(tokenList)

	ast := parser.ParseCode()

	for _, codeString := range ast.CodeStrings {
		parser.Execute(codeString)
	}
}
