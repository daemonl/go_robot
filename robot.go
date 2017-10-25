package robot

import (
	"fmt"
	"strconv"
	"strings"
)

// Robot represents a toy robot positioned on a board in a given dimension
type Robot struct {
	Position  []int64     // Current robot position
	Direction *Direction  // Current heading
	Dimension []Direction // The dimension set to use for the robot.
	Max       []int64     // The maximum coordinates of the robot, i.e. the board size
}

// Direction implements a facing on the board. In 2D space this is North, South
// etc, but could be more complex for further dimensions
type Direction struct {
	Name    string                // Used for reporting and placing the robot
	Turns   map[string]*Direction // When turning `string`, the robot will then be facing `Direction`
	Forward []int64               // Add this to the robot's current position to get the next position after moving forward
}

const (
	North string = "NORTH"
	South string = "SOUTH"
	East  string = "EAST"
	West  string = "WEST"

	Left  string = "LEFT"
	Right string = "RIGHT"
)

// DirectionSet2D implements a two dimentional plane with rotations Left and Right
var DirectionSet2D = []Direction{
	{
		Name:    North,
		Forward: []int64{0, 1},
	},
	{
		Name:    East,
		Forward: []int64{1, 0},
	},
	{
		Name:    South,
		Forward: []int64{0, -1},
	},
	{
		Name:    West,
		Forward: []int64{-1, 0},
	},
}

func init() {
	// Implement Left and Right for 2D, easy as it is a circle.
	for i := 0; i < 4; i++ {
		DirectionSet2D[i].Turns = map[string]*Direction{
			Left:  &DirectionSet2D[(i+3)%4],
			Right: &DirectionSet2D[(i+1)%4],
		}
	}
}

func invalidCommand(tpl string, d ...interface{}) error {
	return fmt.Errorf(tpl, d...)
}

var ErrorWouldFall = fmt.Errorf("Robot would fall")
var ErrorNotPlaced = fmt.Errorf("Robot not placed")

// Report returns a textual representation of the current position and direction
func (r Robot) Report() (string, error) {
	if r.Direction == nil || len(r.Position) == 0 {
		return "", ErrorNotPlaced
	}
	components := make([]string, len(r.Position), len(r.Position))
	for idx, pos := range r.Position {
		components[idx] = strconv.FormatInt(int64(pos), 10)
	}
	return fmt.Sprintf("%s,%s", strings.Join(components, ","), r.Direction.Name), nil
}

// Place sets the position of the robot. It must contain the same number of
// coordinates as the board (2 for 2D)
func (r *Robot) Place(directionName string, coordinates ...int64) error {
	if len(coordinates) != len(r.Max) {
		return invalidCommand("Can only move in one dimension. Please provide %d coordinate values", len(r.Position))
	}

	for idx := range coordinates {
		if coordinates[idx] > r.Max[idx] {
			return ErrorWouldFall
		}
	}

	var newDirection *Direction
	for _, direction := range r.Dimension {
		if direction.Name == directionName {
			newDirection = &direction
			break
		}
	}
	if newDirection == nil {
		return invalidCommand("Invalid direction '%s' for this dimension", directionName)
	}

	r.Position = coordinates
	r.Direction = newDirection

	return nil
}

// Move moves 1 step in the current direction
func (r *Robot) Move() error {
	if len(r.Position) == 0 {
		return ErrorNotPlaced
	}

	newPosition := r.Position
	for idx := range r.Max {
		moveTo := r.Position[idx] + r.Direction.Forward[idx]
		if moveTo < 0 || moveTo > r.Max[idx] {
			return ErrorWouldFall
		}
		newPosition[idx] = moveTo
	}

	r.Position = newPosition
	return nil
}

// Turn turns the robot from the current direction to a new one
func (r *Robot) Turn(turn string) error {
	if len(r.Position) == 0 {
		return ErrorNotPlaced
	}

	newDirection, ok := r.Direction.Turns[turn]
	if !ok {
		return invalidCommand("Invalid turn direction '%s'", turn)
	}

	r.Direction = newDirection
	return nil
}
