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

        #legend {
            display: flex;
            flex-direction: column;
            overflow: auto;
            gap: 10px;
        }

        .legend-item {
            display: flex;
            flex-direction: column;
        }

        .group {
            display: flex;
            flex-direction: row;
        }

        .group-item {
            padding: 10px;
        }
    </style>
</head>
<body>
    <div class="horizontal-layout">
        <div class="vertical-layout">
            <div id="legend"></div>
        </div>
        <div class="vertical-layout">
            <div id="graph"></div>
        </div>
    </div>
    <script>
        const nodes = [{{range .Nodes}}
            {{.Values}},{{end}}
        ];
        
        const edges = [{{range .Edges}}
            {{.Values}},{{end}}
        ];

        {{if .WithLegend}}
        const legend = document.getElementById("legend");
        const nodesByGroup = {}
        const nodesWithoutGroup = []

        for (const node of nodes) {
            if (node.group != null) {
                if (nodesByGroup[node.group] == null) {
                    nodesByGroup[node.group] = { nodes: [] }
                }

                if (node.leader) {
                    nodesByGroup[node.group].leader = node
                    continue
                }

                nodesByGroup[node.group].nodes.push(node)
                continue
            }

            nodesWithoutGroup.push(node)
        }

        for (const group in nodesByGroup) {
            const groupLeader = nodesByGroup[group].leader
            const groupNodes = nodesByGroup[group].nodes

            if (groupLeader == null) {
                for (const node of groupNodes) {
                    const nodeElement = document.createElement("div")
                    nodeElement.className = "legend-item"
                    nodeElement.innerText = node.label
                    legend.appendChild(nodeElement)
                }
                continue
            }

            const groupElement = document.createElement("div")
            groupElement.className = "legend-item"

            const groupHeaderElement = document.createElement("div")
            groupHeaderElement.className = "group"

            const groupHeaderButton = document.createElement("button")
            groupHeaderButton.innerText = ">"

            const leaderElement = document.createElement("div")
            leaderElement.innerText = groupLeader.label

            groupHeaderElement.appendChild(leaderElement)
            groupHeaderElement.appendChild(groupHeaderButton)

            groupElement.appendChild(groupHeaderElement)

            const nodesContainer = document.createElement("div")
            const nodesElements = []

            for (const node of groupNodes) {
                const nodeElement = document.createElement("div")
                nodeElement.className = "group-item"
                nodeElement.innerText = node.label
                nodesElements.push(nodeElement)
                nodesContainer.appendChild(nodeElement)
            }

            groupHeaderButton.onclick = () => {
                for (const node of nodesElements) {
                    if (node.style.display == "none") {
                        node.style.display = "block"
                        continue
                    }

                    if (node.style.display != "none") {
                        node.style.display = "none"
                        continue
                    }
                }
            }

            groupElement.appendChild(nodesContainer)

            legend.appendChild(groupElement)
        }

        for (const node of nodesWithoutGroup) {
            const nodeElement = document.createElement("div")
            nodeElement.className = "legend-item"
            nodeElement.innerText = node.label
            legend.appendChild(nodeElement)
        }
        {{end}}

        const data = { nodes, edges };

        const options = {{.Options}};

        const container = document.getElementById("graph");
        const network = new vis.Network(container, data, options);

        const zoomHandler = function(params) {
            console.log("Zoom:", params);

            const scale = 1 / params.scale;

            const minSize = options.nodes?.scaling?.label?.min || 10;
            const maxSize = options.nodes?.scaling?.label?.max || 30;

            network.setOptions({
                nodes: {
                    font: {
                        size: Math.min(maxSize, Math.max(minSize, minSize * scale))
                    }
                }
            });
        };

        function throttle(func, delay) {
            let lastCall = 0;
            return function(...args) {
                const now = new Date().getTime();
                if (now - lastCall < delay) return;
                lastCall = now;
                return func.apply(this, args);
            };
        };

        network.on("zoom", throttle(zoomHandler, 100));
    </script>
</body>
</html>
`

var (
	DefaultOptions = map[string]any{
		"nodes": map[string]any{
			"shape": "dot",
			"scaling": map[string]any{
				"min": 10,
				"max": 30,
				"label": map[string]any{
					"min": 10,
					"max": 100,
				},
			},
			"font": map[string]any{
				"size": 12,
				"face": "Tahoma",
			},
		},
		"edges": map[string]any{
			"width": 0.15,
			"color": map[string]any{"inherit": "from"},
			"smooth": map[string]any{
				"type": "continuous",
			},
		},
		"physics": map[string]any{
			"solver":        "forceAtlas2Based",
			"stabilization": false,
			"barnesHut": map[string]any{
				"gravitationalConstant": -80000,
				"springConstant":        0.01,
				"springLength":          200,
			},
			"forceAtlas2Based": map[string]any{
				"gravitationalConstant": -2000,
				"springConstant":        0.1,
				"springLength":          200,
			},
		},
		"interaction": map[string]any{
			"tooltipDelay":    200,
			"hideEdgesOnDrag": true,
		},
	}
)

func (container *Graph[T]) HtmlRender(writer io.Writer, options ...map[string]any) error {
	type NodeData struct {
		Values map[string]any
	}

	type EdgeData struct {
		Values map[string]any
	}

	type TemplateData struct {
		Nodes      []NodeData
		Edges      []EdgeData
		Options    map[string]any
		WithLegend bool
	}

	data := TemplateData{
		Nodes:   []NodeData{},
		Edges:   []EdgeData{},
		Options: DefaultOptions,
	}

	for _, optionsMap := range options {
		for key, value := range optionsMap {
			data.Options[key] = value
		}
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
			transitionValues := map[string]any{
				"from":   node.name,
				"to":     transition.node.name,
				"arrows": "to",
			}

			for optionName, optionValue := range transition.options {
				transitionValues[optionName] = optionValue
			}

			data.Edges = append(data.Edges, EdgeData{
				Values: transitionValues,
			})
		}
	}

	tmpl, err := template.New("graph").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, data)
}

func (container *Graph[T]) HtmlRenderToFile(filename string, options ...map[string]any) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return container.HtmlRender(file, options...)
}
