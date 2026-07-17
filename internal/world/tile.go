package world

type TileType int

const (
	TileEmpty TileType = iota
	TileWall
	TileCastle
)

func (t TileType) String() string {
	switch t {
	case TileEmpty:
		return "Empty"
	case TileWall:
		return "Wall"
	case TileCastle:
		return "Castle"
	default:
		return "Unknown"
	}
}
