package container_graph

import (
	"fmt"
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
		html,
        body {
            padding: 0;
            offset: 0;
            display: flex;
            width: 100vw;
            height: 100vh;
        }

        .vertical-layout {
            display: flex;
            flex-grow: 1;
        }

        .horizontal-layout {
            display: flex;
            flex-grow: 1;
        }

        #graph {
            flex-grow: 1;
            border: 1px solid #ddd;
        }
    </style>
</head>
<body>
    <div class="vertical-layout">
        <div class="horizontal-layout">
            <div id="graph"></div>
        </div>
    </div>
    <script>
        const nodes = new vis.DataSet([{{range .Nodes}}
            {{.Values}},{{end}}
        ]);
        
        const edges = new vis.DataSet([{{range .Edges}}
            { from: "{{.From}}", to: "{{.To}}", arrows: "to" },{{end}}
        ]);
        
        const container = document.getElementById("graph");
        const data = { nodes, edges };
        const options = {
            nodes: {
                shape: "box",
                scaling: {
                    min: 10,
                    max: 30,
                },
                font: {
                    size: 12,
                    face: "Tahoma",
                },
            },
            edges: {
                width: 0.15,
                color: { inherit: "from" },
                smooth: {
                    type: "continuous",
                },
            },
            physics: {
                solver: "forceAtlas2Based",
                stabilization: false,
                barnesHut: {
                    gravitationalConstant: -80000,
                    springConstant: 0.01,
                    springLength: 200,
                },
                forceAtlas2Based: {
                    gravitationalConstant: -2000,
                    springConstant: 0.1,
                    springLength: 200,
                },
            },
            interaction: {
                tooltipDelay: 200,
                hideEdgesOnDrag: true,
            },

        };
        new vis.Network(container, data, options);
    </script>
</body>
</html>
`

func (container *Graph[T]) HtmlRender(writer io.Writer) error {
	type NodeData struct {
		Values map[string]any
	}

	type EdgeData struct {
		From, To string
	}

	type TemplateData struct {
		Nodes []NodeData
		Edges []EdgeData
	}

	data := TemplateData{
		Nodes: []NodeData{},
		Edges: []EdgeData{},
	}

	// Собираем все рёбра
	for _, node := range container.nodes {
		values := map[string]any{
			"id":    node.name,
			"label": fmt.Sprintf("%s\n%v", node.name, node.value),
		}

		for optionName, optionValue := range node.options {
			values[optionName] = optionValue
		}

		data.Nodes = append(data.Nodes, NodeData{Values: values})

		for _, transition := range node.transitions {
			data.Edges = append(data.Edges, EdgeData{
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
