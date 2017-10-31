package robot

import (
	"fmt"
	"strconv"
	"strings"
)

// Robot represents a toy robot positioned on a board in a given dimension
type Robot struct {
	Position  []int64    // Current robot position
	Direction *Direction // Current heading
}

type Board struct {
	Dimension []Direction // The dimension set to use for the robot.
	Max       []int64     // The maximum coordinates of the robot, i.e. the board size

	ExpectedRobotCount int // Exactly this many robots are required before commands will be accepted

	robots map[string]*Robot // The robots on the board.
}

func NewBoard(x, y int64, robotCount int) *Board {
	return &Board{
		Dimension:          DirectionSet2D,
		Max:                []int64{x, y},
		robots:             map[string]*Robot{},
		ExpectedRobotCount: robotCount,
	}
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
var ErrorNotPlaced = fmt.Errorf("Not all robots have been placed")
var ErrorCollision = fmt.Errorf("Robot collision")
var ErrorTooManyRobots = fmt.Errorf("Too many robots")

func (board *Board) Report(robotName string) (string, error) {
	if err := board.acceptingMovements(); err != nil {
		return "", err
	}
	robot, ok := board.robots[robotName]
	if !ok {
		return "", ErrorNotPlaced
	}
	return robot.Report()
}

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
func (board *Board) Place(robotName string, directionName string, coordinates ...int64) error {
	if len(coordinates) != len(board.Max) {
		return invalidCommand("Can only move in one dimension. Please provide %d coordinate values", len(board.Max))
	}

	for idx := range coordinates {
		if coordinates[idx] > board.Max[idx] {
			return ErrorWouldFall
		}
	}

	var newDirection *Direction
	for _, direction := range board.Dimension {
		if direction.Name == directionName {
			newDirection = &direction
			break
		}
	}
	if newDirection == nil {
		return invalidCommand("Invalid direction '%s' for this dimension", directionName)
	}

	if err := board.validatePosition(robotName, coordinates); err != nil {
		return err
	}

	if _, ok := board.robots[robotName]; !ok {
		if len(board.robots)+1 > board.ExpectedRobotCount {
			return ErrorTooManyRobots
		}
	}

	board.robots[robotName] = &Robot{
		Position:  coordinates,
		Direction: newDirection,
	}

	return nil
}

// Move moves 1 step in the current direction
func (board *Board) Move(robotName string) error {
	if err := board.acceptingMovements(); err != nil {
		return err
	}
	robot, ok := board.robots[robotName]
	if !ok {
		return ErrorNotPlaced

	}

	newPosition := make([]int64, len(robot.Position), len(robot.Position))
	for idx := range robot.Position {
		newPosition[idx] = robot.Position[idx] + robot.Direction.Forward[idx]
	}

	if err := board.validatePosition(robotName, newPosition); err != nil {
		return err
	}

	robot.Position = newPosition
	return nil
}

func (board Board) acceptingMovements() error {
	if len(board.robots) != board.ExpectedRobotCount {
		return ErrorNotPlaced
	}
	return nil
}

func (board Board) validatePosition(robotName string, newPosition []int64) error {
	for idx := range newPosition {
		if newPosition[idx] < 0 || newPosition[idx] > board.Max[idx] {
			return ErrorWouldFall
		}
	}

robotLoop:
	for otherRobotName, robot := range board.robots {
		if otherRobotName == robotName {
			continue
		}

		for idx := range newPosition {
			if newPosition[idx] != robot.Position[idx] {
				continue robotLoop
			}
		}
		return ErrorCollision
	}

	return nil
}

func (board *Board) Turn(robotName string, direction string) error {
	if err := board.acceptingMovements(); err != nil {
		return err
	}
	robot, ok := board.robots[robotName]
	if !ok {
		return ErrorNotPlaced
	}
	return robot.Turn(direction)
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
