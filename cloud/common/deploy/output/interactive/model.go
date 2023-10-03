package interactive

import (
	"fmt"
	"io"
	"strings"
	"time"

	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/samber/lo"
)

type DeployModel struct {
	pulumiSub chan events.EngineEvent
	sub       chan tea.Msg
	spinner   spinner.Model
	logs      []string
	tree      *Tree[PulumiData]
}

func (m DeployModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		subscribeToChan(m.sub),
		subscribeToChan(m.pulumiSub),
	)
}

// subscribeToChannel - A tea Command that will wait on messages sent to the given channel
func subscribeToChan[T any](sub chan T) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

const MAX_LOG_LENGTH = 5

// Implement io.Writer for simplicity
func (m DeployModel) Write(bytes []byte) (int, error) {
	msg := string(bytes)
	cutMsg, _ := strings.CutSuffix(msg, "\n")

	// This will hook the writer into the tea program lifecycle
	m.sub <- LogMessage{
		Message: cutMsg,
	}

	return len(bytes), nil
}

func (m *DeployModel) handlePulumiEngineEvent(evt events.EngineEvent) {
	// These events are directly tied to a resource
	if evt.DiagnosticEvent != nil {
		// TODO: Handle diagnostic event logging
		node := m.tree.FindNode(evt.DiagnosticEvent.URN)
		if node != nil {
			node.Data.LastMessage = evt.DiagnosticEvent.Message
		}

	} else if evt.ResourcePreEvent != nil {
		// attempt to locate the parent node
		meta := evt.ResourcePreEvent.Metadata.New
		if meta == nil {
			meta = evt.ResourcePreEvent.Metadata.Old
		}

		parentNode := m.tree.FindNode(meta.Parent)
		if parentNode != nil {
			parentNode.AddChild(&Node[PulumiData]{
				Id: evt.ResourcePreEvent.Metadata.URN,
				Data: &PulumiData{
					Urn:    evt.ResourcePreEvent.Metadata.URN,
					Type:   evt.ResourcePreEvent.Metadata.Type,
					Status: PreResourceStates[string(evt.ResourcePreEvent.Metadata.Op)],
				},
			})
		} else {
			m.tree.Root.AddChild(&Node[PulumiData]{
				Id: evt.ResourcePreEvent.Metadata.URN,
				Data: &PulumiData{
					Urn:    evt.ResourcePreEvent.Metadata.URN,
					Type:   evt.ResourcePreEvent.Metadata.Type,
					Status: PreResourceStates[string(evt.ResourcePreEvent.Metadata.Op)],
				},
				Children: []*Node[PulumiData]{},
			})
		}

	} else if evt.ResOutputsEvent != nil {
		// Find the URN in the tree
		node := m.tree.FindNode(evt.ResOutputsEvent.Metadata.URN)
		if node != nil {
			node.Data.Status = SuccessResourceStates[string(evt.ResOutputsEvent.Metadata.Op)]
		}
	} else if evt.ResOpFailedEvent != nil {
		node := m.tree.FindNode(evt.ResOpFailedEvent.Metadata.URN)
		if node != nil {
			node.Data.Status = FailedResourceStates[string(evt.ResOpFailedEvent.Metadata.Op)]
		}
	}
}

func (m DeployModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch t := msg.(type) {
	case events.EngineEvent:
		m.handlePulumiEngineEvent(t)
		return m, subscribeToChan(m.pulumiSub)
	case LogMessage:
		m.logs = append(m.logs, t.Message)
		return m, subscribeToChan(m.sub)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m DeployModel) renderNodeRow(node *Node[PulumiData], depth int, isLast bool, parentLast bool) table.Row {
	linkChar := lo.Ternary(!isLast, "├─", "└─")
	prefixString := lo.Ternary(!parentLast, fmt.Sprintf("│  %s", linkChar), linkChar)
	marginLeft := lo.Ternary(!parentLast, 3*(depth-1), 3*depth)

	statusStyle := StatusStyles[node.Data.Status]
	isPending := lo.Contains(lo.Values(PreResourceStates), node.Data.Status)
	isComplete := lo.Contains(lo.Values(SuccessResourceStates), node.Data.Status)
	isFailed := lo.Contains(lo.Values(FailedResourceStates), node.Data.Status)

	status := statusStyle.Render(MessageResourceStates[node.Data.Status])
	if isPending {
		status = statusStyle.Render(MessageResourceStates[node.Data.Status] + " " + m.spinner.View())
	} else if isComplete {
		status = statusStyle.Render(MessageResourceStates[node.Data.Status] + " " + "👍")
	} else if isFailed {
		status = statusStyle.Render(MessageResourceStates[node.Data.Status] + " " + "👎")
	}

	return table.Row{
		// Name
		lipgloss.NewStyle().MarginLeft(marginLeft).SetString(prefixString).Render(node.Data.Name()),
		// Type
		node.Data.Type,
		// Status
		status,
		// Message
		node.Data.LastMessage,
	}
}

// Render the tree rows
func (m DeployModel) renderNodeRows(depth int, parentLast bool, nodes ...*Node[PulumiData]) []table.Row {
	// render this nods info
	rows := []table.Row{}

	for idx, n := range nodes {
		isLast := idx == len(nodes)-1
		rows = append(rows, m.renderNodeRow(n, depth, isLast, parentLast))

		if len(n.Children) > 0 {
			rows = append(rows, m.renderNodeRows(depth+1, isLast, n.Children...)...)
		}
	}

	return rows
}

func (m DeployModel) View() string {
	rows := m.renderNodeRows(0, true, m.tree.Root)

	// Calculate dynamic column widths
	// colWidths := lo.Reduce(rows, func(agg []int, row table.Row, idx int) []int {
	// 	return []int{
	// 		lo.Max[int]([]int{agg[0], len(row[0])}),
	// 		lo.Max[int]([]int{agg[1], len(row[1])}),
	// 		lo.Max[int]([]int{agg[2], len(row[2])}),
	// 		lo.Max[int]([]int{agg[3], len(row[3])}),
	// 	}
	// }, []int{0, 0, 0, 0})

	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Type", Width: 20},
		{Title: "Status", Width: 20},
		{Title: "Logs", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Background(lipgloss.Color("28")).
		BorderBottom(true).
		Bold(false)

	t.SetStyles(s)

	return fmt.Sprintf("\n%s\n", t.View())
}

// TODO: Take the message we get from up
func NewInteractiveOutput(sub chan tea.Msg, pulumiSub chan events.EngineEvent, output io.Writer) *tea.Program {
	// Return the new tea program without running it
	return tea.NewProgram(DeployModel{
		pulumiSub: pulumiSub,
		sub:       sub,
		logs:      make([]string, 0),
		spinner:   spinner.New(spinner.WithSpinner(spinner.Ellipsis)),
		tree: &Tree[PulumiData]{
			Root: &Node[PulumiData]{
				Id: "root",
				Data: &PulumiData{
					Urn:    "root",
					Type:   "project",
					Status: ResourceStatus_Unchanged,
				},
				Children: []*Node[PulumiData]{},
			},
		},
	}, tea.WithOutput(output))
}

var HeresHoping = spinner.Spinner{
	Frames: []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"},
	FPS:    time.Second / 24, //nolint:gomnd
}

func NewOutputModel(sub chan tea.Msg, pulumiSub chan events.EngineEvent) DeployModel {
	// FIXME: Set this according to the connected output preferences
	os.Setenv("CLICOLOR_FORCE", "1")

	return DeployModel{
		pulumiSub: pulumiSub,
		sub:       sub,
		logs:      make([]string, 0),
		spinner:   spinner.New(spinner.WithSpinner(spinner.Dot)),
		tree: &Tree[PulumiData]{
			Root: &Node[PulumiData]{
				Id: "root",
				Data: &PulumiData{
					Urn:    "project",
					Type:   "",
					Status: ResourceStatus_None,
				},
				Children: []*Node[PulumiData]{},
			},
		},
	}
}
