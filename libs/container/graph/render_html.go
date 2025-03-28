package container_graph

import (
	"html/template"
	"io"
	"os"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Graph Visualization</title>
    <script src="https://cdn.jsdelivr.net/npm/vis-network@9.1.2/dist/vis-network.min.js"></script>
    <style>
        #graph { width: 100%; height: 600px; border: 1px solid #ddd; }
    </style>
</head>
<body>
    <div id="graph"></div>
    <script>
        const nodes = new vis.DataSet([{{range .Nodes}}
            { id: "{{.Name}}", label: "{{.Name}}\n{{.Value}}", shape: "circle" },{{end}}
        ]);
        
        const edges = new vis.DataSet([{{range .Edges}}
            { from: "{{.From}}", to: "{{.To}}", arrows: "to" },{{end}}
        ]);
        
        const container = document.getElementById("graph");
        const data = { nodes, edges };
        const options = {
            layout: { hierarchical: { direction: "LR" } },
            physics: { hierarchicalRepulsion: { nodeDistance: 120 } }
		};
        new vis.Network(container, data, options);
    </script>
</body>
</html>
`

func (container *Graph[T]) HtmlRender(writer io.Writer) error {
	type Edge struct {
		From, To string
	}
	type TemplateData struct {
		Nodes []*Node[T]
		Edges []Edge
	}

	data := TemplateData{
		Nodes: container.nodes,
	}

	// Собираем все рёбра
	for _, node := range container.nodes {
		for _, transition := range node.transitions {
			data.Edges = append(data.Edges, Edge{
				From: node.name,
				To:   transition.name,
			})
		}
	}

	tmpl, err := template.New("graph").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

func (container *Graph[T]) HtmlRenderToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return container.HtmlRender(file)
}
