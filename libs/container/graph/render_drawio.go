package container_graph

import (
	"encoding/xml"
	"fmt"
	"math"
	"os"

	"github.com/necroin/golibs/utils"
	forcedirected "github.com/necroin/golibs/utils/force-directed"
)

// DrawIOConfig содержит настройки для экспорта в draw.io
type DrawIOConfig struct {
	NodeWidth  int
	NodeHeight int
	SpacingX   int
	SpacingY   int
}

// DefaultDrawIOConfig возвращает дефолтные настройки
func DefaultDrawIOConfig() *DrawIOConfig {
	return &DrawIOConfig{
		NodeWidth:  120,
		NodeHeight: 60,
		SpacingX:   200,
		SpacingY:   100,
	}
}

type mxFile struct {
	XMLName    xml.Name `xml:"mxfile"`
	Host       string   `xml:"host,attr"`
	Modified   string   `xml:"modified,attr"`
	Compressed bool     `xml:"compressed,attr"`
	Pages      []mxPage `xml:"diagram"`
}

type mxPage struct {
	XMLName xml.Name `xml:"diagram"`
	Name    string   `xml:"name,attr"`
	Id      string   `xml:"id,attr"`
	Graph   mxGraph  `xml:"mxGraphModel"`
}

type mxGraph struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Size    int      `xml:"gridSize,attr"`
	Cells   []mxCell `xml:"root>mxCell"`
}

type mxCell struct {
	ID       string      `xml:"id,attr,omitempty"`
	Value    string      `xml:"value,attr,omitempty"`
	Style    string      `xml:"style,attr,omitempty"`
	Parent   string      `xml:"parent,attr,omitempty"`
	Source   string      `xml:"source,attr,omitempty"`
	Target   string      `xml:"target,attr,omitempty"`
	Edge     string      `xml:"edge,attr,omitempty"`
	Vertex   string      `xml:"vertex,attr,omitempty"`
	Geometry *mxGeometry `xml:"mxGeometry,omitempty"`
}

type mxGeometry struct {
	X        int    `xml:"x,attr,omitempty"`
	Y        int    `xml:"y,attr,omitempty"`
	Width    int    `xml:"width,attr,omitempty"`
	Height   int    `xml:"height,attr,omitempty"`
	Relative int    `xml:"relative,attr,omitempty"`
	As       string `xml:"as,attr,omitempty"`
}

func (container *Graph[T]) ExportToDrawIO(filename string, config *DrawIOConfig) error {
	if config == nil {
		config = DefaultDrawIOConfig()
	}

	// Инициализация позиций
	positions := map[string]*utils.Vector2D{}
	for _, node := range container.nodes {
		positions[node.name] = &utils.Vector2D{}
	}

	forcedirected.InitializePositions(
		container.NodesNames(),
		positions,
	)

	// Оптимизация расположения
	forcedirected.OptimizeLayout(
		container.NodesNames(),
		container.NodesTransitionsNames(),
		positions,
		forcedirected.ForceDirectedParams{
			Repulsion:     100,
			Stiffness:     1.2,
			CoolingFactor: 0.98,
			Iterations:    300,
		},
	)

	// Нормализация координат
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64

	for _, pos := range positions {
		if pos.X < minX {
			minX = pos.X
		}
		if pos.Y < minY {
			minY = pos.Y
		}
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}

	scale := min(
		float64(config.SpacingX*len(container.nodes))/(maxX-minX+100),
		float64(config.SpacingY*len(container.nodes))/(maxY-minY+100),
	)

	mx := mxFile{
		Host:     "Electron",
		Modified: "2023-01-01T00:00:00Z",
		Pages: []mxPage{{
			Name: "Page-1",
			Id:   "page1",
			Graph: mxGraph{
				Size: 10,
				Cells: []mxCell{
					// Root cell
					{ID: "0"},
					// Default parent cell
					{ID: "1", Parent: "0"},
				},
			},
		}},
	}

	x, y := 0, 0
	nodeId := 2 // Начинаем с 2, так как 0 и 1 уже заняты
	nodeMap := make(map[string]string)

	// Добавляем узлы
	for _, node := range container.nodes {
		cellID := fmt.Sprintf("%d", nodeId)
		nodeMap[node.name] = cellID

		mx.Pages[0].Graph.Cells = append(mx.Pages[0].Graph.Cells, mxCell{
			ID:     cellID,
			Value:  fmt.Sprintf("%s\n%v", node.name, node.value),
			Style:  "rounded=1;whiteSpace=wrap;html=1;fillColor=#ffffff;strokeColor=#000000;",
			Parent: "1",
			Vertex: "1",
			Geometry: &mxGeometry{
				X:      int((positions[node.name].X - minX) * scale),
				Y:      int((positions[node.name].Y - minY) * scale),
				Width:  config.NodeWidth,
				Height: config.NodeHeight,
				As:     "geometry",
			},
		})

		x += config.SpacingX
		if x > config.SpacingX*3 {
			x = 0
			y += config.SpacingY
		}
		nodeId++
	}

	// Добавляем связи
	for _, node := range container.nodes {
		for _, transition := range node.transitions {
			sourceID, ok := nodeMap[node.name]
			if !ok {
				continue
			}

			targetID, ok := nodeMap[transition.name]
			if !ok {
				continue
			}

			mx.Pages[0].Graph.Cells = append(mx.Pages[0].Graph.Cells, mxCell{
				ID:     fmt.Sprintf("%d", nodeId),
				Source: sourceID,
				Target: targetID,
				Edge:   "1",
				Parent: "1",
				Style:  "rounded=0;html=1;jettySize=auto;orthogonalLoop=1;",
				Geometry: &mxGeometry{
					As:       "geometry",
					Width:    50,
					Height:   50,
					Relative: 1,
				},
			})
			nodeId++
		}
	}

	output, err := xml.MarshalIndent(mx, "", "  ")
	if err != nil {
		return fmt.Errorf("xml marshal error: %w", err)
	}

	header := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	xmlData := append([]byte(header), output...)

	return os.WriteFile(filename, xmlData, 0644)
}
