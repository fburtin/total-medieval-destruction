package world

type Grid struct {
	columns int
	rows    int
	tiles   [][]TileType
}

func NewGrid(columns, rows int) *Grid {
	tiles := make([][]TileType, rows)

	for row := 0; row < rows; row++ {
		tiles[row] = make([]TileType, columns)
	}

	return &Grid{
		columns: columns,
		rows:    rows,
		tiles:   tiles,
	}
}

func (g *Grid) Columns() int {
	return g.columns
}

func (g *Grid) Rows() int {
	return g.rows
}

func (g *Grid) IsInside(column, row int) bool {
	return column >= 0 &&
		column < g.columns &&
		row >= 0 &&
		row < g.rows
}

func (g *Grid) TileAt(column, row int) TileType {
	if !g.IsInside(column, row) {
		return TileEmpty
	}

	return g.tiles[row][column]
}

func (g *Grid) SetTile(column, row int, tile TileType) {
	if !g.IsInside(column, row) {
		return
	}

	g.tiles[row][column] = tile
}

func (g *Grid) ClearTile(column, row int) {
	g.SetTile(column, row, TileEmpty)
}

func (g *Grid) ToggleWall(column, row int) {
	if !g.IsInside(column, row) {
		return
	}

	currentTile := g.TileAt(column, row)

	switch currentTile {
	case TileEmpty:
		g.SetTile(column, row, TileWall)

	case TileWall:
		g.SetTile(column, row, TileEmpty)

	case TileCastle, TileWater:
		return
	}
}

func (g *Grid) SetCastle(column, row int) {
	if !g.IsInside(column, row) {
		return
	}

	if g.TileAt(column, row) != TileEmpty {
		return
	}

	g.removeExistingCastle()
	g.SetTile(column, row, TileCastle)
}

func (g *Grid) IsBuildable(column, row int) bool {
	if !g.IsInside(column, row) {
		return false
	}

	tile := g.TileAt(column, row)

	return tile == TileEmpty || tile == TileWall
}

func (g *Grid) FindCastle() (column int, row int, found bool) {
	for currentRow := 0; currentRow < g.rows; currentRow++ {
		for currentColumn := 0; currentColumn < g.columns; currentColumn++ {
			if g.tiles[currentRow][currentColumn] == TileCastle {
				return currentColumn, currentRow, true
			}
		}
	}

	return 0, 0, false
}

func (g *Grid) removeExistingCastle() {
	column, row, found := g.FindCastle()

	if found {
		g.SetTile(column, row, TileEmpty)
	}
}
