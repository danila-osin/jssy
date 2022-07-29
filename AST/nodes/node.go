package nodes

import (
	tok "github.com/strange-cat-fe/jssy/token"
)

type Node struct {
	Token *tok.Token
	LeftNode *Node
	RightNode *Node
}

func NewNode(token *tok.Token, leftNode *Node, rightNode *Node) *Node {
	return &Node{Token: token, LeftNode: leftNode, RightNode: rightNode}
}