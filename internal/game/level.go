package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/fburtin/total-medieval-destruction/internal/entities"
	"github.com/fburtin/total-medieval-destruction/internal/world"
)

const (
	riverWidth              = 4
	riverGenerationAttempts = 100
)

type riverPosition struct {
	column int
	row    int
}

type riverPoint struct {
	column float64
	row    float64
}

func (g *Game) initializeLevel() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	g.generateBalancedRiver(rng)
	g.spawnEnemyShips(rng)
}

func (g *Game) generateBalancedRiver(rng *rand.Rand) {
	var bestRiver [][]bool
	bestDifference := math.MaxInt

	for attempt := 0; attempt < riverGenerationAttempts; attempt++ {
		river := g.createSingleRiver(rng)

		firstArea, secondArea := g.calculateTwoLandAreas(river)

		if firstArea == 0 || secondArea == 0 {
			continue
		}

		difference := absInt(firstArea - secondArea)

		if difference < bestDifference {
			bestDifference = difference
			bestRiver = river
		}
	}

	if bestRiver != nil {
		g.applyRiver(bestRiver)
	}
}

func (g *Game) createSingleRiver(rng *rand.Rand) [][]bool {
	rows := g.grid.Rows()
	columns := g.grid.Columns()

	river := make([][]bool, rows)

	for row := range river {
		river[row] = make([]bool, columns)
	}

	start, end := g.randomOppositeBorders(rng)

	deltaColumn := end.column - start.column
	deltaRow := end.row - start.row

	length := math.Hypot(deltaColumn, deltaRow)

	if length == 0 {
		return river
	}

	perpendicularColumn := -deltaRow / length
	perpendicularRow := deltaColumn / length

	steps := int(length * 4)

	currentCurve := 0.0
	targetCurve := 0.0

	for step := 0; step <= steps; step++ {
		progress := float64(step) / float64(steps)

		if step%12 == 0 {
			targetCurve = rng.Float64()*6 - 3
		}

		currentCurve += (targetCurve - currentCurve) * 0.08

		centerColumn :=
			start.column +
				deltaColumn*progress +
				perpendicularColumn*currentCurve

		centerRow :=
			start.row +
				deltaRow*progress +
				perpendicularRow*currentCurve

		drawRiverDisk(
			river,
			centerColumn,
			centerRow,
			float64(riverWidth)/2,
		)
	}

	return river
}

func (g *Game) randomOppositeBorders(
	rng *rand.Rand,
) (riverPoint, riverPoint) {
	columns := g.grid.Columns()
	rows := g.grid.Rows()

	switch rng.Intn(2) {
	case 0:
		// Top to bottom.
		return riverPoint{
				column: float64(rng.Intn(columns)),
				row:    0,
			},
			riverPoint{
				column: float64(rng.Intn(columns)),
				row:    float64(rows - 1),
			}

	default:
		// Left to right.
		return riverPoint{
				column: 0,
				row:    float64(rng.Intn(rows)),
			},
			riverPoint{
				column: float64(columns - 1),
				row:    float64(rng.Intn(rows)),
			}
	}
}

func drawRiverDisk(
	river [][]bool,
	centerColumn float64,
	centerRow float64,
	radius float64,
) {
	minColumn := int(math.Floor(centerColumn - radius))
	maxColumn := int(math.Ceil(centerColumn + radius))

	minRow := int(math.Floor(centerRow - radius))
	maxRow := int(math.Ceil(centerRow + radius))

	for row := minRow; row <= maxRow; row++ {
		if row < 0 || row >= len(river) {
			continue
		}

		for column := minColumn; column <= maxColumn; column++ {
			if column < 0 || column >= len(river[row]) {
				continue
			}

			distance := math.Hypot(
				float64(column)-centerColumn,
				float64(row)-centerRow,
			)

			if distance <= radius {
				river[row][column] = true
			}
		}
	}
}

func (g *Game) calculateTwoLandAreas(
	river [][]bool,
) (int, int) {
	rows := g.grid.Rows()
	columns := g.grid.Columns()

	visited := make([][]bool, rows)

	for row := range visited {
		visited[row] = make([]bool, columns)
	}

	areas := make([]int, 0, 2)

	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			if river[row][column] || visited[row][column] {
				continue
			}

			area := floodFillArea(
				river,
				visited,
				column,
				row,
			)

			areas = append(areas, area)
		}
	}

	if len(areas) != 2 {
		return 0, 0
	}

	return areas[0], areas[1]
}

func floodFillArea(
	river [][]bool,
	visited [][]bool,
	startColumn int,
	startRow int,
) int {
	type position struct {
		column int
		row    int
	}

	queue := []position{
		{
			column: startColumn,
			row:    startRow,
		},
	}

	visited[startRow][startColumn] = true
	area := 0

	directions := []position{
		{column: 1},
		{column: -1},
		{row: 1},
		{row: -1},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		area++

		for _, direction := range directions {
			nextColumn := current.column + direction.column
			nextRow := current.row + direction.row

			if nextRow < 0 ||
				nextRow >= len(river) ||
				nextColumn < 0 ||
				nextColumn >= len(river[nextRow]) {
				continue
			}

			if river[nextRow][nextColumn] ||
				visited[nextRow][nextColumn] {
				continue
			}

			visited[nextRow][nextColumn] = true

			queue = append(
				queue,
				position{
					column: nextColumn,
					row:    nextRow,
				},
			)
		}
	}

	return area
}

func (g *Game) applyRiver(river [][]bool) {
	for row := 0; row < g.grid.Rows(); row++ {
		for column := 0; column < g.grid.Columns(); column++ {
			if river[row][column] {
				g.grid.SetTile(
					column,
					row,
					world.TileWater,
				)
			}
		}
	}
}

func (g *Game) spawnEnemyShips(rng *rand.Rand) {
	const (
		minShips = 1
		maxShips = 5

		shipWidthInTiles  = 3
		shipHeightInTiles = 3
	)

	shipCount := minShips + rng.Intn(maxShips-minShips+1)

	validPositions := make([]riverPosition, 0)

	for row := 0; row <= g.grid.Rows()-shipHeightInTiles; row++ {
		for column := 0; column <= g.grid.Columns()-shipWidthInTiles; column++ {
			if !g.isWaterRectangle(
				column,
				row,
				shipWidthInTiles,
				shipHeightInTiles,
			) {
				continue
			}

			validPositions = append(
				validPositions,
				riverPosition{
					column: column,
					row:    row,
				},
			)
		}
	}

	rng.Shuffle(
		len(validPositions),
		func(first, second int) {
			validPositions[first], validPositions[second] =
				validPositions[second], validPositions[first]
		},
	)

	g.enemyShips = make([]*entities.Ship, 0, shipCount)

	for _, position := range validPositions {
		if len(g.enemyShips) >= shipCount {
			break
		}

		if g.isTooCloseToExistingShip(position.column, position.row) {
			continue
		}

		x := float64(
			gridOffsetX + position.column*tileSize,
		)

		y := float64(
			gridOffsetY + position.row*tileSize,
		)

		g.enemyShips = append(
			g.enemyShips,
			entities.NewShip(x, y),
		)
	}
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

func (g *Game) isWaterRectangle(
	startColumn int,
	startRow int,
	width int,
	height int,
) bool {
	for row := startRow; row < startRow+height; row++ {
		for column := startColumn; column < startColumn+width; column++ {
			if !g.grid.IsInside(column, row) {
				return false
			}

			if g.grid.TileAt(column, row) != world.TileWater {
				return false
			}
		}
	}

	return true
}

func (g *Game) isTooCloseToExistingShip(
	column int,
	row int,
) bool {
	const minimumDistanceInTiles = 5

	for _, ship := range g.enemyShips {
		shipColumn := int(
			(ship.X - float64(gridOffsetX)) /
				float64(tileSize),
		)

		shipRow := int(
			(ship.Y - float64(gridOffsetY)) /
				float64(tileSize),
		)

		columnDistance := absInt(column - shipColumn)
		rowDistance := absInt(row - shipRow)

		if columnDistance < minimumDistanceInTiles &&
			rowDistance < minimumDistanceInTiles {
			return true
		}
	}

	return false
}
