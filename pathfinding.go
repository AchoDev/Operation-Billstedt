package main

import (
	"container/heap"
	"fmt"
	"math"
)

// Point represents a point in the grid
type Point struct {
	X, Y int
}

// Node represents a node in the priority queue
type Node struct {
	Point    Point
	Priority float64
	Index    int
}

// PriorityQueue implements a priority queue for A* algorithm
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.Index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.Index = -1
	*pq = old[0 : n-1]
	return node
}

// Heuristic calculates the Manhattan distance between two points
func Heuristic(a, b Point) float64 {
	return math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
}

// Neighbors returns the neighboring points of a given point
func Neighbors(p Point, grid [][]int) []Point {
	directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	var neighbors []Point

	for _, d := range directions {
		np := Point{p.X + d.X, p.Y + d.Y}
		if np.X >= 0 && np.X < len(grid) && np.Y >= 0 && np.Y < len(grid[0]) && grid[np.X][np.Y] == 0 {
			neighbors = append(neighbors, np)
		}
	}

	return neighbors
}

// AStar implements the A* algorithm
func AStar(start, goal Point, grid [][]int) []Point {
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Node{Point: start, Priority: 0})

	cameFrom := make(map[Point]Point)
	gScore := make(map[Point]float64)
	gScore[start] = 0

	fScore := make(map[Point]float64)
	fScore[start] = Heuristic(start, goal)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Node).Point

		if current == goal {
			var path []Point
			for current != start {
				path = append([]Point{current}, path...)
				current = cameFrom[current]
			}
			path = append([]Point{start}, path...)
			return path
		}

		for _, neighbor := range Neighbors(current, grid) {
			tentativeGScore := gScore[current] + 1

			if g, ok := gScore[neighbor]; !ok || tentativeGScore < g {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + Heuristic(neighbor, goal)
				heap.Push(pq, &Node{Point: neighbor, Priority: fScore[neighbor]})
			}
		}
	}

	return nil // No path found
}

func runPathfindingAlgorithm(start, end Transform, colliders []*Collider, gridSize int, worldSize Vector2) []Vector2 {
	// grid size = size of grid tile in pixels
	// world size = size of world in pixels

	if isInsideCollider(start.GetPosition(), colliders) {
		fmt.Println("start is inside a collider")
		return []Vector2{}
	}

	if isInsideCollider(end.GetPosition(), colliders) {
		fmt.Println("end is inside a collider")
		return []Vector2{}
	}

	minWorldX := -2000.0
	minWorldY := -2000.0
	maxWorldX := worldSize.x + math.Abs(minWorldX)
	maxWorldY := worldSize.y + math.Abs(minWorldY)

	// Calculate grid dimensions and offsets
	gridWidth := int((maxWorldX - minWorldX) / float64(gridSize))
	gridHeight := int((maxWorldY - minWorldY) / float64(gridSize))

	// Initialize grid
	grid := make([][]int, gridWidth)
	for x := 0; x < gridWidth; x++ {
		grid[x] = make([]int, gridHeight)
	}

	// Mark colliders in the grid
	for _, collider := range colliders {
		x := int((collider.transform.x - minWorldX) / float64(gridSize))
		y := int((collider.transform.y - minWorldY) / float64(gridSize))
		width := int(collider.transform.width / float64(gridSize))
		height := int(collider.transform.height / float64(gridSize))

		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				gridX := x + i
				gridY := y + j
				if gridX >= 0 && gridX < gridWidth && gridY >= 0 && gridY < gridHeight {
					grid[gridX][gridY] = 1
				}
			}
		}
	}

	// Convert start and end positions to grid coordinates (with offset)
	startPoint := Point{
		X: int((start.x - minWorldX) / float64(gridSize)),
		Y: int((start.y - minWorldY) / float64(gridSize)),
	}
	endPoint := Point{
		X: int((end.x - minWorldX) / float64(gridSize)),
		Y: int((end.y - minWorldY) / float64(gridSize)),
	}

	path := AStar(startPoint, endPoint, grid)

	vectorPath := []Vector2{}
	for _, point := range path {
		// Convert grid coordinates back to world coordinates
		worldX := float64(point.X)*float64(gridSize) + minWorldX
		worldY := float64(point.Y)*float64(gridSize) + minWorldY
		vectorPath = append(vectorPath, Vector2{
			x: worldX,
			y: worldY,
		})
	}

	if len(path) == 0 {
		fmt.Println("pathfinding returns no path")
		return []Vector2{}
	}

	return vectorPath

	// nextPoint := path[0]
	// return Vector2{
	//     x: float64(nextPoint.X * gridSize),
	//     y: float64(nextPoint.Y * gridSize),
	// }
}

func isInsideCollider(position Vector2, colliders []*Collider) bool {
    for _, collider := range colliders {
        if position.x >= collider.transform.x &&
            position.x <= collider.transform.x+collider.transform.width &&
            position.y >= collider.transform.y &&
            position.y <= collider.transform.y+collider.transform.height {
            return true
        }
    }
    return false
}
