package forcedirected

import (
	"math"

	"github.com/necroin/golibs/utils"
)

const minDistance = 50

type ForceDirectedParams struct {
	Repulsion     float64 // Сила отталкивания
	Stiffness     float64 // Жесткость связей
	CoolingFactor float64 // Коэффициент охлаждения
	Iterations    int
}

// Рассчитывает силы между узлами
func CalculateForces(nodes []string, neighbors map[string][]string, positions map[string]*utils.Vector2D, forces map[string]*utils.Vector2D, params ForceDirectedParams) {
	// Отталкивание между всеми узлами
	for _, node1 := range nodes {
		for _, node2 := range nodes {
			if node1 == node2 {
				continue
			}

			vector := utils.Vector2D{
				X: positions[node1].X - positions[node2].X,
				Y: positions[node1].Y - positions[node2].Y,
			}

			distance := max(vector.Distance(), 0.01)

			if distance < minDistance {
				adjust := (minDistance - distance) / 2
				forces[node1].X += adjust * vector.X / distance
				forces[node1].Y += adjust * vector.Y / distance
				continue
			}

			force := math.Log(distance * distance)
			forces[node1].X += force * vector.X / distance
			forces[node1].Y += force * vector.Y / distance
		}
	}

	// Притяжение по связям
	for _, node := range nodes {
		for _, neighbor := range neighbors[node] {
			vector := utils.Vector2D{
				X: positions[neighbor].X - positions[node].X,
				Y: positions[neighbor].Y - positions[node].Y,
			}
			distance := vector.Distance()
			force := params.Stiffness * distance

			forces[node].X += force * vector.X / distance
			forces[node].Y += force * vector.Y / distance
		}
	}
}

// Оптимизирует расположение узлов
func OptimizeLayout(nodes []string, neighbors map[string][]string, positions map[string]*utils.Vector2D, params ForceDirectedParams) {
	forces := map[string]*utils.Vector2D{}
	velocities := map[string]*utils.Vector2D{}
	cooling := 1.0

	for name := range positions {
		velocities[name] = &utils.Vector2D{}
	}

	for range params.Iterations {
		// Сброс сил
		for name := range positions {
			forces[name] = &utils.Vector2D{}
		}

		params.Repulsion *= cooling
		params.Stiffness *= cooling

		// Расчет сил
		CalculateForces(nodes, neighbors, positions, forces, params)

		// Обновление позиций
		for name := range positions {
			velocities[name].X += 0.8*velocities[name].X + 0.2*forces[name].X
			velocities[name].Y += 0.8*velocities[name].Y + 0.2*forces[name].Y

			positions[name].X += velocities[name].X * cooling
			positions[name].Y += velocities[name].Y * cooling
		}

		cooling *= params.CoolingFactor
	}
}

// Вместо случайного размещения
func InitializePositions(nodes []string, positions map[string]*utils.Vector2D) {
	// Размещение узлов по кругу
	centerX, centerY := 500.0, 500.0
	radius := 300.0
	angleStep := 2 * math.Pi / float64(len(nodes))

	for i, node := range nodes {
		angle := float64(i) * angleStep
		positions[node] = &utils.Vector2D{
			X: centerX + radius*math.Cos(angle),
			Y: centerY + radius*math.Sin(angle),
		}
	}
}
