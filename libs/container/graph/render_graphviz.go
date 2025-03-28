package container_graph

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func (container *Graph[T]) VisualizeDOT() string {
	var builder strings.Builder

	builder.WriteString("digraph G {\n")
	builder.WriteString("  rankdir=LR;\n") // Ориентация слева направо
	builder.WriteString("  node [shape=circle];\n")

	// Добавляем все узлы
	for _, node := range container.nodes {
		builder.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\\n%v\"];\n", node.name, node.name, node.value))
	}

	// Добавляем все переходы
	for _, node := range container.nodes {
		for _, transition := range node.transitions {
			builder.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", node.name, transition.name))
		}
	}

	builder.WriteString("}\n")
	return builder.String()
}

func (container *Graph[T]) GraphvizRender(ctx context.Context, writer io.Writer, format graphviz.Format, shape cgraph.Shape) error {
	graphviz, _ := graphviz.New(ctx)
	defer graphviz.Close()

	graph, err := graphviz.Graph()
	if err != nil {
		return fmt.Errorf("failed to create graph: %w", err)
	}
	defer graph.Close()

	graph.SetRankDir(cgraph.LRRank)

	cgraphNodes := make(map[string]*cgraph.Node)
	for _, node := range container.nodes {
		cgraphNode, err := graph.CreateNodeByName(node.name)
		if err != nil {
			return fmt.Errorf("failed to create node %s: %w", node.name, err)
		}
		cgraphNode.SetLabel(fmt.Sprintf("%s\n%v", node.name, node.value))
		cgraphNode.SetShape(shape)
		cgraphNode.SetStyle(cgraph.FilledNodeStyle)
		cgraphNodes[node.name] = cgraphNode
	}

	for _, node := range container.nodes {
		for _, transition := range node.transitions {
			_, err := graph.CreateEdgeByName("", cgraphNodes[node.name], cgraphNodes[transition.name])
			if err != nil {
				return fmt.Errorf("failed to create edge %s->%s: %w", node.name, transition.name, err)
			}
		}
	}

	if err := graphviz.Render(ctx, graph, format, writer); err != nil {
		return fmt.Errorf("failed to render graph: %w", err)
	}

	return nil
}

func (container *Graph[T]) GraphvizRenderToFile(ctx context.Context, filename string, format graphviz.Format, shape cgraph.Shape) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()
	return container.GraphvizRender(ctx, file, format, shape)
}
