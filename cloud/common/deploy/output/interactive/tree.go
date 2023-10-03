package interactive

type Tree[T any] struct {
	Root *Node[T]
}

// Find node from root
func (t *Tree[T]) FindNode(id string) *Node[T] {
	return t.findNode(t.Root, id)
}

// Recursive implementation
func (t *Tree[T]) findNode(node *Node[T], id string) *Node[T] {
	if node.Id == id {
		// Return if we have a match
		return node
	}

	// otherwise walk the children recursively
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			if foundNode := t.findNode(child, id); foundNode != nil {
				return foundNode
			}
		}
	}

	// No matches
	return nil
}

type Node[T any] struct {
	Id       string
	Data     *T
	Children []*Node[T]
}

func (n *Node[T]) AddChild(node *Node[T]) {
	if node.Children == nil {
		node.Children = make([]*Node[T], 0)
	}

	n.Children = append(n.Children, node)
}
