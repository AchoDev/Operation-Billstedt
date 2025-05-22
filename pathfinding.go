package main

import (
	"container/heap"
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
    directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {1, 1}, {-1, -1}, {1, -1}, {-1, 1}}
    var neighbors []Point
    for _, d := range directions {
        np := Point{p.X + d.X, p.Y + d.Y}
        if np.X >= 0 && np.X < len(grid) && np.Y >= 0 && np.Y < len(grid[0]) && grid[np.X][np.Y] == 0 {
            neighbors = append(neighbors, np)
        }
    }
    return neighbors
}

// AStar implements the A* algorithm with cycle prevention and max iteration cap
func AStar(start, goal Point, grid [][]int) []Point {
    const maxIterations = 20000

    pq := &PriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, &Node{Point: start, Priority: 0})

    cameFrom := make(map[Point]Point)
    gScore := make(map[Point]float64)
    gScore[start] = 0
    fScore := make(map[Point]float64)
    fScore[start] = Heuristic(start, goal)
    visited := make(map[Point]bool)

    iterationCount := 0

    for pq.Len() > 0 {
        if iterationCount > maxIterations {
            return nil
        }
        iterationCount++

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

        if visited[current] {
            continue
        }
        visited[current] = true

        for _, neighbor := range Neighbors(current, grid) {
            if visited[neighbor] {
                continue
            }

            // Add the cost of the cell to the gScore
            cellCost := float64(grid[neighbor.X][neighbor.Y])
            if cellCost == 1 { // Wall
                continue
            }

            tentativeGScore := gScore[current] + 1 + cellCost
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

// isInsideCollider checks if a given position is inside any collider
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

func moveOutOfCollider(position Point, grid [][]int) Point {

    gridWidth := len(grid)
    gridHeight := len(grid[0])

    if grid[position.X][position.Y] > 0 {
        directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {1, 1}, {-1, -1}, {1, -1}, {-1, 1}}
        for _, d := range directions {
            adjustedPoint := Point{position.X + d.X, position.Y + d.Y}
            if adjustedPoint.X >= 0 && adjustedPoint.X < gridWidth &&
            adjustedPoint.Y >= 0 && adjustedPoint.Y < gridHeight &&
            grid[adjustedPoint.X][adjustedPoint.Y] != 1 {
                position = adjustedPoint
            break
            }
        }

        for _, d := range directions {
            adjustedPoint := Point{position.X + d.X, position.Y + d.Y}
            if adjustedPoint.X >= 0 && adjustedPoint.X < gridWidth &&
            adjustedPoint.Y >= 0 && adjustedPoint.Y < gridHeight &&
            grid[adjustedPoint.X][adjustedPoint.Y] == 0 {
                position = adjustedPoint
            break
            }
        }
    }

    return position
}

// runPathfindingAlgorithm initializes the grid, marks colliders, and runs A*
func runPathfindingAlgorithm(start, end Transform, colliders []*Collider, gridSize int, worldSize Vector2) []Vector2 {
    minWorldX, minWorldY := -14000.0, -14000.0
    maxWorldX := worldSize.x + math.Abs(minWorldX)
    maxWorldY := worldSize.y + math.Abs(minWorldY)

    gridWidth := int((maxWorldX - minWorldX) / float64(gridSize))
    gridHeight := int((maxWorldY - minWorldY) / float64(gridSize))

    grid := make([][]int, gridWidth)
    for x := 0; x < gridWidth; x++ {
        grid[x] = make([]int, gridHeight)
    }

    for _, collider := range colliders {
        halfW, halfH := collider.transform.width/2, collider.transform.height/2
        startX := int((collider.transform.x - halfW - minWorldX) / float64(gridSize))
        startY := int((collider.transform.y - halfH - minWorldY) / float64(gridSize))
        endX := int((collider.transform.x + halfW - minWorldX) / float64(gridSize))
        endY := int((collider.transform.y + halfH - minWorldY) / float64(gridSize))

        bufferCells := int(math.Ceil(1.0 / float64(gridSize)))
		
        for gridX := startX; gridX <= endX; gridX++ {
            for gridY := startY; gridY <= endY; gridY++ {
                if gridX >= 0 && gridX < gridWidth && gridY >= 0 && gridY < gridHeight {
                    grid[gridX][gridY] = 1
                }
            }
        }

        for gridX := startX - bufferCells; gridX <= endX + bufferCells; gridX++ {
            for gridY := startY - bufferCells; gridY <= endY + bufferCells; gridY++ {
                if gridX >= 0 && gridX < gridWidth && gridY >= 0 && gridY < gridHeight {
                    if grid[gridX][gridY] == 0 || grid[gridX][gridY] > 10 { 
                        grid[gridX][gridY] = 10
                    }
                }
            }
        }
    }

    // Check if the target grid position is inside a collider
    endPoint := Point{
        X: int((end.x - minWorldX) / float64(gridSize)),
        Y: int((end.y - minWorldY) / float64(gridSize)),
    }
    // Check if the start grid position is inside a collider
    startPoint := Point{
        X: int((start.x - minWorldX) / float64(gridSize)),
        Y: int((start.y - minWorldY) / float64(gridSize)),
    }

    startPoint = moveOutOfCollider(startPoint, grid)
    endPoint = moveOutOfCollider(endPoint, grid)

    path := AStar(startPoint, endPoint, grid)

    var vectorPath []Vector2
    for _, point := range path {
        vectorPath = append(vectorPath, Vector2{
            x: float64(point.X)*float64(gridSize) + minWorldX,
            y: float64(point.Y)*float64(gridSize) + minWorldY,
        })
    }

    return vectorPath
}