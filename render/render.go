package render

import "github.com/hajimehoshi/ebiten/v2"

type Direction uint8

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
	UPLEFT
	UPRIGHT
	DOWNRIGHT
	DOWNLEFT
)

var (
	// Usage: x to subtract when something is following cusor
	XOffset float64
	// Usage: y to subtract when something is following cusor
	YOffset float64
	// Usage: DirectionMap[Direction].X * math.Abs(renderable.Velocity.X)
	DirectionMap = map[Direction]Vec2{
		UP:        Vec2{X: 0, Y: -1},
		DOWN:      Vec2{X: 0, Y: 1},
		LEFT:      Vec2{X: -1, Y: 0},
		RIGHT:     Vec2{X: 1, Y: 0},
		UPLEFT:    Vec2{X: -1, Y: 1},
		UPRIGHT:   Vec2{X: 1, Y: -1},
		DOWNRIGHT: Vec2{X: 1, Y: 1},
		DOWNLEFT:  Vec2{X: -1, Y: 1},
	}
)

// Renderable implements methods that allow
// for the Draw() section of the game loop to
// render anywhere
type Renderable interface {
	Update()               // When the sprite needs to change position between frames
	Sprite() *ebiten.Image // the image to draw
	Position() Vec2        // where the image should go
}
