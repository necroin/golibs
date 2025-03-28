package container_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/goccy/go-graphviz"
	container_graph "github.com/necroin/golibs/libs/container/graph"
	"github.com/necroin/golibs/utils"
	"github.com/necroin/golibs/utils/generator"
)

func TestGraph_Render(t *testing.T) {
	generator := generator.New(false)

	newNode := func() *container_graph.Node[int] {
		value := generator.Next()
		return container_graph.NewNode(
			fmt.Sprintf("Node %d", value),
			value,
		)
	}

	nodes := []*container_graph.Node[int]{
		newNode(),
	}

	graph := container_graph.New(nodes)

	for range 100 {
		node := newNode()
		graph.AddNodeItem(node)
		nodes = append(nodes, node)
	}

	for _, node := range graph.Nodes() {
		randomNode := utils.GetRandomFrom(graph.Nodes()...)
		node.AddTransition(randomNode)
	}

	if err := graph.TopologicalSort(); err != nil {
		t.Log(err)
	}

	if err := utils.SaveToFile("graph.dot", []byte(graph.VisualizeDOT())); err != nil {
		t.Log(err)
	}

	if err := graph.HtmlRenderToFile("graph.html"); err != nil {
		t.Log(err)
	}

	if err := graph.GraphvizRenderToFile(context.Background(), "graph.svg", graphviz.SVG, graphviz.CircleShape); err != nil {
		t.Log(err)
	}
}
