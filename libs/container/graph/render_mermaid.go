package container_graph

import (
	"fmt"
	"strings"
)

type MermaidOptions struct {
	Direction     string
	MakeSubgraphs bool
}

type MermaidOption func(options *MermaidOptions)

func (container *Graph[T]) VisualizeMermaid(opts ...MermaidOption) string {
	options := &MermaidOptions{
		Direction:     "LR",
		MakeSubgraphs: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	var builder strings.Builder

	builder.WriteString("```mermaid\n")
	builder.WriteString(fmt.Sprintf("graph %s\n", options.Direction))

	if options.MakeSubgraphs {
		groups := map[string][]string{}

		for _, node := range container.nodes {
			groupOption := node.Options()["group"]
			if groupOption != nil {
				group := fmt.Sprintf("%s", groupOption)
				groups[group] = append(groups[group], node.Name())
			}
		}

		for group, nodes := range groups {
			builder.WriteString(fmt.Sprintf("\tsubgraph %s\n", group))

			for _, node := range nodes {
				builder.WriteString(fmt.Sprintf("\t\t%s\n", node))
			}

			builder.WriteString("\tend\n")
		}

	}

	for _, node := range container.nodes {
		for _, transition := range node.Transitions() {
			builder.WriteString(fmt.Sprintf("\t%s --> %s\n", node.Name(), transition.Node().Name()))
		}
	}

	builder.WriteString("```\n")
	return builder.String()
}

func WithMermaidDicrection(direction string) MermaidOption {
	return func(options *MermaidOptions) {
		options.Direction = direction
	}
}

func WithMermaidGroup() MermaidOption {
	return func(options *MermaidOptions) {
		options.MakeSubgraphs = true
	}
}
