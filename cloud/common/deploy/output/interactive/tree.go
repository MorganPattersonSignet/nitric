package interactive

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Tree[T any] struct {
	Root *Node[T]
}

// Find node from root
func (t *Tree[T]) FindNode(urn string) *Node[T] {
	return t.findNode(t.Root, urn)
}

// Recursive implementation
func (t *Tree[T]) findNode(node *Node[T], urn string) *Node[T] {
	if node.Id == urn {
		// Return if we have a match
		return node
	}

	// otherwise walk the children recursively
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			if foundNode := t.findNode(child, urn); foundNode != nil {
				return foundNode
			}
		}
	}

	// No matches
	return nil
}

type ResourceStatus int

const (
	ResourceStatus_Creating = iota
	ResourceStatus_Updating
	ResourceStatus_Deleting
	ResourceStatus_Created
	ResourceStatus_Deleted
	ResourceStatus_Updated
	ResourceStatus_Failed_Create
	ResourceStatus_Failed_Delete
	ResourceStatus_Failed_Update
	ResourceStatus_Unchanged
)

var PreResourceStates = map[string]ResourceStatus{
	"create": ResourceStatus_Creating,
	"delete": ResourceStatus_Deleting,
	"same":   ResourceStatus_Unchanged,
	"update": ResourceStatus_Updating,
}

var SuccessResourceStates = map[string]ResourceStatus{
	"create": ResourceStatus_Created,
	"delete": ResourceStatus_Deleted,
	"same":   ResourceStatus_Unchanged,
	"update": ResourceStatus_Updated,
}

var FailedResourceStates = map[string]ResourceStatus{
	"create": ResourceStatus_Failed_Create,
	"delete": ResourceStatus_Failed_Delete,
	"update": ResourceStatus_Failed_Update,
}

var MessageResourceStates = map[ResourceStatus]string{
	ResourceStatus_Creating:      "creating",
	ResourceStatus_Updating:      "updating",
	ResourceStatus_Deleting:      "deleting",
	ResourceStatus_Created:       "created",
	ResourceStatus_Deleted:       "deleted",
	ResourceStatus_Updated:       "updated",
	ResourceStatus_Failed_Create: "create failed",
	ResourceStatus_Failed_Delete: "delete failed",
	ResourceStatus_Failed_Update: "updated failed",
	ResourceStatus_Unchanged:     "unchanged",
}

// TODO: Use TUI standard colors when lib available
var StatusStyles = map[ResourceStatus]lipgloss.Style{
	// Unchanged State
	ResourceStatus_Unchanged: lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC")),
	// Pre states
	ResourceStatus_Creating: lipgloss.NewStyle().Foreground(lipgloss.Color("#2ecc71")),
	ResourceStatus_Updating: lipgloss.NewStyle().Foreground(lipgloss.Color("#f1c40f")),
	ResourceStatus_Deleting: lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c")),
	// Post states
	ResourceStatus_Created: lipgloss.NewStyle().Foreground(lipgloss.Color("#27ae60")),
	ResourceStatus_Updated: lipgloss.NewStyle().Foreground(lipgloss.Color("#f39c12")),
	ResourceStatus_Deleted: lipgloss.NewStyle().Foreground(lipgloss.Color("#c0392b")),

	// Failed states
	ResourceStatus_Failed_Create: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
	ResourceStatus_Failed_Delete: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
	ResourceStatus_Failed_Update: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
}

type PulumiData struct {
	Urn string
	// Name   string
	Type   string
	Status ResourceStatus
}

func (pd PulumiData) Name() string {
	urnParts := strings.Split(pd.Urn, "::")

	return urnParts[len(urnParts)-1]
}

type IdentifiableNode interface {
	Id() string
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
