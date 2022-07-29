package nodes

type StatementsNode struct {
	CodeStrings []*Node
}

func (sn *StatementsNode) AddNode(node *Node) {
	sn.CodeStrings = append(sn.CodeStrings, node)
}
