package world

type gridPosition struct {
	column int
	row    int
}

func (g *Grid) IsCastleEnclosed() bool {
	castleColumn, castleRow, found := g.FindCastle()

	if !found {
		return false
	}

	visited := make([][]bool, g.rows)

	for row := 0; row < g.rows; row++ {
		visited[row] = make([]bool, g.columns)
	}

	queue := []gridPosition{
		{
			column: castleColumn,
			row:    castleRow,
		},
	}

	visited[castleRow][castleColumn] = true

	directions := []gridPosition{
		{column: 1, row: 0},
		{column: -1, row: 0},
		{column: 0, row: 1},
		{column: 0, row: -1},
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if g.isBoundary(current.column, current.row) {
			return false
		}

		for _, direction := range directions {
			nextColumn := current.column + direction.column
			nextRow := current.row + direction.row

			if !g.IsInside(nextColumn, nextRow) {
				continue
			}

			if visited[nextRow][nextColumn] {
				continue
			}

			if g.TileAt(nextColumn, nextRow) == TileWall {
				continue
			}

			visited[nextRow][nextColumn] = true

			queue = append(queue, gridPosition{
				column: nextColumn,
				row:    nextRow,
			})
		}
	}

	return true
}

func (g *Grid) isBoundary(column, row int) bool {
	return column == 0 ||
		column == g.columns-1 ||
		row == 0 ||
		row == g.rows-1
}
