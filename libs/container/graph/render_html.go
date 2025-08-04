package container_graph

import (
	"fmt"
	"html/template"
	"io"
	"os"
)

const htmlZoomHandler = `
        const zoomHandler = function(params) {
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

        network.on("zoom", zoomHandler);
`

const htmlHighlightHandler = `
        var highlightActive = false;
        var allNodes = nodesDataset.get({ returnType: "Object" });

        function neighbourhoodHighlight(params) {
            // if something is selected:
            if (params.nodes.length > 0) {
                highlightActive = true;
                var i, j;
                var selectedNode = params.nodes[0];
                var degrees = 2;

                // mark all nodes as hard to read.
                for (var nodeId in allNodes) {
                    allNodes[nodeId].color = "rgba(200,200,200,0.5)";
                    if (allNodes[nodeId].hiddenLabel === undefined) {
                        allNodes[nodeId].hiddenLabel = allNodes[nodeId].label;
                        allNodes[nodeId].label = undefined;
                    }
                }
                var connectedNodes = network.getConnectedNodes(selectedNode);
                var allConnectedNodes = [];

                // get the second degree nodes
                for (i = 1; i < degrees; i++) {
                    for (j = 0; j < connectedNodes.length; j++) {
                        allConnectedNodes = allConnectedNodes.concat(
                            network.getConnectedNodes(connectedNodes[j])
                        );
                    }
                }

                // all second degree nodes get a different color and their label back
                for (i = 0; i < allConnectedNodes.length; i++) {
                    allNodes[allConnectedNodes[i]].color = "rgba(150,150,150,0.75)";
                    if (allNodes[allConnectedNodes[i]].hiddenLabel !== undefined) {
                        allNodes[allConnectedNodes[i]].label =
                            allNodes[allConnectedNodes[i]].hiddenLabel;
                        allNodes[allConnectedNodes[i]].hiddenLabel = undefined;
                    }
                }

                // all first degree nodes get their own color and their label back
                for (i = 0; i < connectedNodes.length; i++) {
                    allNodes[connectedNodes[i]].color = undefined;
                    if (allNodes[connectedNodes[i]].hiddenLabel !== undefined) {
                        allNodes[connectedNodes[i]].label =
                            allNodes[connectedNodes[i]].hiddenLabel;
                        allNodes[connectedNodes[i]].hiddenLabel = undefined;
                    }
                }

                // the main node gets its own color and its label back.
                allNodes[selectedNode].color = undefined;
                if (allNodes[selectedNode].hiddenLabel !== undefined) {
                    allNodes[selectedNode].label = allNodes[selectedNode].hiddenLabel;
                    allNodes[selectedNode].hiddenLabel = undefined;
                }
            } else if (highlightActive === true) {
                // reset all nodes
                for (var nodeId in allNodes) {
                    allNodes[nodeId].color = undefined;
                    if (allNodes[nodeId].hiddenLabel !== undefined) {
                        allNodes[nodeId].label = allNodes[nodeId].hiddenLabel;
                        allNodes[nodeId].hiddenLabel = undefined;
                    }
                }
                highlightActive = false;
            }

            // transform the object into an array
            var updateArray = [];
            for (nodeId in allNodes) {
                if (allNodes.hasOwnProperty(nodeId)) {
                    updateArray.push(allNodes[nodeId]);
                }
            }
            nodesDataset.update(updateArray);
        }

        network.on("click", neighbourhoodHighlight);
`

const htmlClustering = `
        function groupNodesBySimilarEdges(nodes, edges) {
            const nodeConnections = {};

            // Собираем связи для каждого узла
            nodes.forEach(node => {
                nodeConnections[node.id] = new Set();
            });
        
            edges.forEach(edge => {
                nodeConnections[edge.from].add(edge.to);
                nodeConnections[edge.to].add(edge.from);
            });
        
            // Группируем узлы с одинаковыми связями
            const clusters = {};
            nodes.forEach(node => {
                const key = [...nodeConnections[node.id]].sort().join(',');
                if (!clusters[key]) clusters[key] = [];
                clusters[key].push(node);
            });
        
            return Object.values(clusters);
        }

        const clusters = groupNodesBySimilarEdges(nodes, edges);
        console.log(clusters);

        for (const cluster of clusters) {
            if (clusters.length >= 1) {
                const ids = [];
                for (const node of cluster){
                    ids.push(node.id)
                }

                network.cluster({
                    joinCondition: (childOptions) => ids.includes(childOptions.id),
                    clusterNodeProperties: {
                        label: cluster[0].group,
                        group: cluster[0].group,
                        shape: "dot",
                    },
                })
            }
        }
`

const htmlLegend = `
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
`

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

        const nodesDataset = new vis.DataSet(nodes);
        const edgesDataset = new vis.DataSet(edges);

        const data = { nodes: nodesDataset, edges: edgesDataset };

        const options = {{.Options}};

        const container = document.getElementById("graph");
        const network = new vis.Network(container, data, options);

        console.log(network);

        {{.Legend}}
        {{.Clustering}}
        {{.ZoomHandler}}
        {{.HighlightHandler}}
        
        {{if .PhysicsStopDelay}}
        setTimeout(() => { network.setOptions({ physics: { enabled: false } }) }, {{.PhysicsStopDelay}});
        {{end}}
    </script>
</body>
</html>
`

type HtmlNodeData struct {
	Values map[string]any
}

type HtmlEdgeData struct {
	Values map[string]any
}

type HtmlTemplateData struct {
	Nodes   []HtmlNodeData
	Edges   []HtmlEdgeData
	Options map[string]any

	Legend           template.JS
	Clustering       template.JS
	ZoomHandler      template.JS
	HighlightHandler template.JS
	PhysicsStopDelay int64
}

type HtmlOption func(data *HtmlTemplateData)

var (
	DefaultHtmlOptions = map[string]any{
		"nodes": map[string]any{
			"shape": "dot",
			"scaling": map[string]any{
				"min": 10,
				"max": 30,
			},
			"font": map[string]any{
				"size": 12,
				"face": "Tahoma",
			},
		},
		"edges": map[string]any{
			"width": 0.05,
			"color": map[string]any{"inherit": "from"},
			"smooth": map[string]any{
				"type": "continuous",
			},
		},
		"physics": map[string]any{
			"solver": "forceAtlas2Based",
			"stabilization": map[string]any{
				"enabled":    true,
				"iterations": 100000,
			},
			"barnesHut": map[string]any{
				"gravitationalConstant": -2000,
				"springConstant":        0.01,
				"springLength":          200,
			},
			"forceAtlas2Based": map[string]any{
				"gravitationalConstant": -2000,
				"springConstant":        0.1,
				// "springLength":          200,

				"theta":          0.5,
				"centralGravity": 0.01,
				"springLength":   100,
				"damping":        0.4,
				"avoidOverlap":   1,
			},
		},
		"interaction": map[string]any{
			"tooltipDelay":    200,
			"hideEdgesOnDrag": true,
			"hideEdgesOnZoom": true,
		},
	}
)

func (container *Graph[T]) HtmlRender(writer io.Writer, options ...HtmlOption) error {
	data := &HtmlTemplateData{
		Nodes:   []HtmlNodeData{},
		Edges:   []HtmlEdgeData{},
		Options: DefaultHtmlOptions,
	}

	for _, option := range options {
		option(data)
	}

	// Собираем все рёбра
	for _, node := range container.nodes {
		values := map[string]any{
			"id":    node.Name(),
			"label": fmt.Sprintf("%s\n%v", node.Name(), node.Value()),
		}

		for optionName, optionValue := range node.Options() {
			values[optionName] = optionValue
		}

		data.Nodes = append(data.Nodes, HtmlNodeData{Values: values})

		for _, transition := range node.Transitions() {
			transitionValues := map[string]any{
				"from":   node.Name(),
				"to":     transition.Node().Name(),
				"arrows": "to",
			}

			for optionName, optionValue := range transition.Options() {
				transitionValues[optionName] = optionValue
			}

			data.Edges = append(data.Edges, HtmlEdgeData{
				Values: transitionValues,
			})
		}
	}

	tmpl, err := template.New("graph").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(writer, *data)
}

func (container *Graph[T]) HtmlRenderToFile(filename string, options ...HtmlOption) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return container.HtmlRender(file, options...)
}

func WithNetworkOptions(options map[string]any) HtmlOption {
	return func(data *HtmlTemplateData) {
		for key, value := range options {
			data.Options[key] = value
		}
	}
}

func WithLegend() HtmlOption {
	return func(data *HtmlTemplateData) {
		data.Legend = template.JS(htmlLegend)
	}
}

func WithClustering() HtmlOption {
	return func(data *HtmlTemplateData) {
		data.Clustering = template.JS(htmlClustering)
	}
}

func WithZoom() HtmlOption {
	return func(data *HtmlTemplateData) {
		data.ZoomHandler = template.JS(htmlZoomHandler)
	}
}

func WithHighlight() HtmlOption {
	return func(data *HtmlTemplateData) {
		data.HighlightHandler = template.JS(htmlHighlightHandler)
	}
}

func WithPhysicsStopDelay(delay int64) HtmlOption {
	return func(data *HtmlTemplateData) {
		data.PhysicsStopDelay = delay
	}
}
