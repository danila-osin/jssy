package parser

import (
	"fmt"
	"github.com/strange-cat-fe/jssy/AST/nodes"
	tok "github.com/strange-cat-fe/jssy/token"
	"log"
	"strconv"
)

type Parser struct {
	Tokens   []tok.Token
	Position int
	Scope    map[string]interface{}
}

func NewParser(tokens []tok.Token) *Parser {
	return &Parser{Tokens: tokens, Scope: make(map[string]interface{})}
}

func (p *Parser) matchAndGoToNext(expectedTypes []string) *tok.Token {

	tokenTypes := tok.GetTypes(expectedTypes)

	if p.Position < len(p.Tokens) {
		currentToken := p.Tokens[p.Position]

		for _, tokenType := range tokenTypes {
			if tokenType.Name == currentToken.Type.Name {
				p.Position += 1
				return &currentToken
			}
		}

		return nil
	}

	return nil
}

func (p *Parser) require(expectedTypes ...string) *tok.Token {
	token := p.matchAndGoToNext(expectedTypes)

	if token == nil {
		log.Fatalf("Expected token %s on position %d", expectedTypes[0], p.Position)
	}

	return token
}

func (p *Parser) parseExpression() *nodes.Node {
	leftPar := p.matchAndGoToNext([]string{"leftPar"})

	if leftPar != nil {
		node := p.parseFormula()
		p.require("rightPar")
		return node
	} else {
		node := p.parseNumberOrVariable()
		return node
	}
}

func (p *Parser) parseFormula() *nodes.Node {
	leftNode := p.parseExpression()

	operator := p.matchAndGoToNext([]string{"plus", "minus"})

	if operator == nil {
		return leftNode
		//log.Fatalf("") // TODO
	}

	rightNode := p.parseExpression()

	expressionNode := nodes.NewNode(operator, leftNode, rightNode)

	return expressionNode
}

func (p *Parser) parseNumberOrVariable() *nodes.Node {
	number := p.matchAndGoToNext([]string{"number"})
	if number != nil {
		return nodes.NewNode(number, nil, nil)
	}

	variable := p.matchAndGoToNext([]string{"variable"})
	if variable != nil {
		return nodes.NewNode(variable, nil, nil)
	}

	panic("Expected number or variable parse numOrVar")
}

func (p *Parser) parseCodeString() *nodes.Node {
	number := p.matchAndGoToNext([]string{"number"})
	if number != nil {
		panic("Expected variable on position: " + string(rune(p.Position)))
	}

	variable := p.matchAndGoToNext([]string{"variable"})
	if variable != nil {
		variableNode := nodes.NewNode(variable, nil, nil)

		assignOperator := p.matchAndGoToNext([]string{"assign"})

		rightFormulaNode := p.parseFormula()

		if assignOperator == nil {
			log.Fatalf("Assign operation is expected after varOrNum on position: %d", p.Position)
		}

		codeStringNode := nodes.NewNode(assignOperator, variableNode, rightFormulaNode)

		return codeStringNode
	}

	varDecl := p.matchAndGoToNext([]string{"varDeclaration"})
	if varDecl != nil {
		varDeclNode := nodes.NewNode(varDecl, nil, nil)

		assignOperator := p.matchAndGoToNext([]string{"assign"})

		rightFormulaNode := p.parseFormula()

		if assignOperator == nil {
			log.Fatalf("Assign operation is expected after varOrNum on position: %d", p.Position)
		}

		codeStringNode := nodes.NewNode(assignOperator, varDeclNode, rightFormulaNode)

		return codeStringNode
	}

	printFn := p.matchAndGoToNext([]string{"printFn"})
	if printFn != nil {
		rightFormulaNode := p.parseFormula()

		codeStringNode := nodes.NewNode(printFn, nil, rightFormulaNode)

		return codeStringNode
	}

	panic("Expected variable on position: " + string(rune(p.Position)))
}

func (p *Parser) ParseCode() nodes.StatementsNode {
	root := new(nodes.StatementsNode)

	for p.Position < len(p.Tokens) {
		codeStringNode := p.parseCodeString()

		p.require("semicolon")

		root.AddNode(codeStringNode)
	}

	return *root
}

func (p *Parser) Execute(node *nodes.Node) interface{} {
	nodeType := node.Token.Type.Name
	nodeText := node.Token.Text

	switch nodeType {
	case "number":
		number, _ := strconv.ParseFloat(nodeText, 64)
		return number
	case "printFn":
		fmt.Println(p.Execute(node.RightNode))
		return 0
	case "plus":
		leftOperand := p.Execute(node.LeftNode)
		rightOperand := p.Execute(node.RightNode)

		return leftOperand.(float64) + rightOperand.(float64)
	case "minus":
		leftOperand := p.Execute(node.LeftNode)
		rightOperand := p.Execute(node.RightNode)

		return leftOperand.(float64) - rightOperand.(float64)
	case "assign":
		expressionResult := p.Execute(node.RightNode)
		variableNode := node.LeftNode

		if variableNode.Token.Type.Name == "varDeclaration" {
			_, ok := p.Scope[variableNode.Token.Text]
			if ok {
				panic("Variable " + variableNode.Token.Text + " already declared!")
			} else {
				p.Scope[variableNode.Token.Text[4:]] = expressionResult
			}
		} else {
			_, ok := p.Scope[variableNode.Token.Text]
			if ok {
				p.Scope[variableNode.Token.Text] = expressionResult
			} else {
				panic("Variable " + variableNode.Token.Text + " isn't declared!")
			}
		}

		return expressionResult
	case "variable":
		varFromScope, ok := p.Scope[nodeText]
		if ok {
			return varFromScope
		} else {
			panic("Undeclared variable: " + nodeText)
		}

	default:
		panic("Unexpected Error")
	}
}
