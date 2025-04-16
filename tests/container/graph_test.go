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

func AllRender[T any](t *testing.T, graph *container_graph.Graph[T]) {
	if err := utils.SaveToFile("graph.dot", []byte(graph.VisualizeDOT())); err != nil {
		t.Log(err)
	}

	if err := graph.HtmlRenderToFile("graph.html"); err != nil {
		t.Log(err)
	}

	if err := graph.GraphvizRenderToFile(context.Background(), "graph.svg", graphviz.SVG, graphviz.CircleShape); err != nil {
		t.Log(err)
	}

	if err := graph.ExportToDrawIO("graph.drawio", nil); err != nil {
		t.Log(err)
	}
}

func TestGraph_Random_Render(t *testing.T) {
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

	graph := container_graph.New(nodes...)

	for range 100 {
		node := newNode()
		graph.AddNodeItem(node)
		nodes = append(nodes, node)
	}

	for _, node := range graph.Nodes() {
		randomNode := utils.GetRandomFrom(graph.Nodes()...)
		node.AddTransition(randomNode)
	}

	// if err := graph.TopologicalSort(); err != nil {
	// 	t.Log(err)
	// }

	AllRender(t, graph)
}

func TestGraph_Preset_Render(t *testing.T) {
	nodes := []*container_graph.Node[string]{
		container_graph.NewNode("Server_1", "Addr 1").SetOption("group", 1),
		container_graph.NewNode("Server_1-App_1", "Port 1").SetOption("group", 1),
		container_graph.NewNode("Server_1-App_2", "Port 2").SetOption("group", 1),
		container_graph.NewNode("Server_1-App_3", "Port 3").SetOption("group", 1),
		container_graph.NewNode("Server_1-App_4", "Port 4").SetOption("group", 1),

		container_graph.NewNode("Server_2", "Addr 2").SetOption("group", 2),
		container_graph.NewNode("Server_2-App_1", "Port 1").SetOption("group", 2),
		container_graph.NewNode("Server_2-App_2", "Port 2").SetOption("group", 2),
		container_graph.NewNode("Server_2-App_3", "Port 3").SetOption("group", 2),
		container_graph.NewNode("Server_2-App_4", "Port 4").SetOption("group", 2),

		container_graph.NewNode("Server_3", "Addr 2").SetOption("group", 3),
		container_graph.NewNode("Server_3-App_1", "Port 1").SetOption("group", 3),
		container_graph.NewNode("Server_3-App_2", "Port 2").SetOption("group", 3),
		container_graph.NewNode("Server_3-App_3", "Port 3").SetOption("group", 3),
		container_graph.NewNode("Server_3-App_4", "Port 4").SetOption("group", 3),
	}

	graph := container_graph.New(nodes...)

	graph.AddTransitionUndirected("Server_1", "Server_1-App_1")
	graph.AddTransitionUndirected("Server_1", "Server_1-App_2")
	graph.AddTransitionUndirected("Server_1", "Server_1-App_3")
	graph.AddTransitionUndirected("Server_1", "Server_1-App_4")

	graph.AddTransitionUndirected("Server_2", "Server_2-App_1")
	graph.AddTransitionUndirected("Server_2", "Server_2-App_2")
	graph.AddTransitionUndirected("Server_2", "Server_2-App_3")
	graph.AddTransitionUndirected("Server_2", "Server_2-App_4")

	graph.AddTransitionUndirected("Server_3", "Server_3-App_1")
	graph.AddTransitionUndirected("Server_3", "Server_3-App_2")
	graph.AddTransitionUndirected("Server_3", "Server_3-App_3")
	graph.AddTransitionUndirected("Server_3", "Server_3-App_4")

	graph.AddTransition("Server_1", "Server_2")
	graph.AddTransition("Server_3", "Server_2")

	AllRender(t, graph)
}

func TestGraph_Preset_Generate_Render_HTML(t *testing.T) {
	serversCount := 5
	appsCount := 5

	servers := []*container_graph.Node[string]{}
	apps := []*container_graph.Node[string]{}

	for serverId := 0; serverId < serversCount; serverId++ {
		server := container_graph.NewNode(fmt.Sprintf("Server_%d", serverId), fmt.Sprintf("Addr %d", serverId))
		server.SetOption("group", serverId)
		server.SetOption("leader", true)
		servers = append(servers, server)

		for appId := 0; appId < appsCount; appId++ {
			app := container_graph.NewNode(fmt.Sprintf("Server_%d_App_%d", serverId, appId), fmt.Sprintf("Port %d", appId))
			app.SetOption("group", serverId)
			apps = append(apps, app)
		}
	}

	nodes := append(servers, apps...)
	graph := container_graph.New(nodes...)

	for serverId := 0; serverId < serversCount; serverId++ {
		for appId := 0; appId < appsCount; appId++ {
			graph.AddTransition(
				fmt.Sprintf("Server_%d_App_%d", serverId, appId),
				fmt.Sprintf("Server_%d", serverId),
			)
		}
	}

	for _, server1 := range servers {
		for _, server2 := range servers {
			graph.AddTransitionUndirected(
				server1.Name(),
				server2.Name(),
			)
		}
	}

	if err := graph.HtmlRenderToFile("graph.html"); err != nil {
		t.Log(err)
	}
}
