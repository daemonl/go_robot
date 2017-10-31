package robot

import (
	"fmt"
	"testing"
)

func (b Board) assetPosition(robotName string, direction string, coords ...int64) error {
	r, ok := b.robots[robotName]
	if !ok {
		return fmt.Errorf("Is not placed (nil)")
	}

	if r.Direction == nil {
		return fmt.Errorf("Is not placed")
	}

	if r.Direction.Name != direction {
		return fmt.Errorf("Expected direction %s, got %s", direction, r.Direction.Name)
	}

	if len(coords) != len(r.Position) {
		return fmt.Errorf("Robot is in the wrong dimension. Expected %v, got %v", coords, r.Position)
	}

	for idx, expect := range coords {
		if r.Position[idx] != expect {

			return fmt.Errorf("Robot is in the wrong position. Expected %v, got %v", coords, r.Position)
		}
	}

	return nil
}

func TestMulti(t *testing.T) {
	board := NewBoard(1, 1, 2)

	robotA := "foo"
	robotB := "bar"

	if err := board.Place(robotA, North, 0, 0); err != nil {
		t.Fatalf(err.Error())
	}

	if _, err := board.Report(robotB); err == nil {
		t.Error("Expected error")
	} else if err != ErrorNotPlaced {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

	if err := board.Place(robotB, North, 0, 0); err == nil {
		t.Error("Expected error")
	} else if err != ErrorCollision {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

	if err := board.Place(robotB, South, 0, 1); err != nil {
		t.Error(err)
	}

	if err := board.Move(robotB); err == nil {
		t.Error("Expected error")
	} else if err != ErrorCollision {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

	if err := board.Place("baz", North, 1, 1); err == nil {
		t.Error("Expected error")
	} else if err != ErrorTooManyRobots {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

}

func TestPlace(t *testing.T) {
	board := NewBoard(1, 1, 1)

	robotName := "foo"
	if _, err := board.Report(robotName); err == nil {
		t.Error("Expected error")
	} else if err != ErrorNotPlaced {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

	// Wrong number of coordinates
	if err := board.Place(robotName, South); err == nil {
		t.Fatal("Expected error")
	}
	if err := board.Place(robotName, South, 0); err == nil {
		t.Fatal("Expected error")
	}
	if err := board.Place(robotName, South, 0, 0, 0); err == nil {
		t.Fatal("Expected error")
	}

	// Made up direction
	if err := board.Place(robotName, "UP", 0, 0); err == nil {
		t.Fatal("Expected error")
	}

	// Fall
	if err := board.Place(robotName, South, 2, 0); err == nil {
		t.Fatal("Expected error")
	}
	if err := board.Place(robotName, South, 0, 2); err == nil {
		t.Fatal("Expected error")
	}

	if err := board.Place(robotName, North, 0, 0); err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	robot, ok := board.robots[robotName]
	if !ok {
		t.Fatalf("Robot did not exist")
	}

	reported, err := board.Report(robotName)
	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}
	if reported != "0,0,NORTH" {
		t.Errorf("Bad Report: %s", reported)
	}

	if robot.Direction.Name != North {
		t.Errorf("Expected North, got %s", robot.Direction.Name)
	}

	if len(robot.Position) != 2 ||
		robot.Position[0] != 0 ||
		robot.Position[1] != 0 {
		t.Errorf("Wrong Position: %v", robot.Position)
	}
}

func TestMove(t *testing.T) {
	board := NewBoard(2, 2, 1)
	robotName := "foo"
	if err := board.Move(robotName); err == nil {
		t.Error("Expected error")
	} else if err != ErrorNotPlaced {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

	// On a 3x3 grid, Move one step from the center in each direction, check
	// position, then try to fall

	for direction, position := range map[string][]int64{
		North: {1, 2},
		East:  {2, 1},
		South: {1, 0},
		West:  {0, 1},
	} {

		if err := board.Place(robotName, direction, 1, 1); err != nil {
			t.Fatalf("Unexpected Error: %s", err.Error())
		}
		if err := board.assetPosition(robotName, direction, 1, 1); err != nil {
			t.Fatalf("Testing %s: %s", direction, err.Error())
		}
		if err := board.Move(robotName); err != nil {
			t.Fatalf("Testing %s: %s", direction, err.Error())
		}
		if err := board.assetPosition(robotName, direction, position...); err != nil {
			t.Fatalf("Testing %s: %s", direction, err.Error())
		}

		// Try to move once more, would fall
		if err := board.Move(robotName); err == nil {
			t.Error("Expected error")
		} else if err != ErrorWouldFall {
			t.Errorf("Wrong error: %s\n", err.Error())
		}
	}

}

func TestTurn(t *testing.T) {

	expectDirections := map[string][]string{
		North: {West, East},
		East:  {North, South},
		South: {East, West},
		West:  {South, North},
	}
	for _, direction := range DirectionSet2D {
		expect := expectDirections[direction.Name]
		if expect[0] != direction.Turns[Left].Name ||
			expect[1] != direction.Turns[Right].Name {
			t.Errorf("From %s, L: %s, R: %s", direction.Name, direction.Turns[Left].Name, direction.Turns[Right].Name)
		}
	}

	board := NewBoard(0, 0, 1)
	robotName := "FOO"
	if err := board.Turn(robotName, Left); err == nil {
		t.Error("Expected error")
	} else if err != ErrorNotPlaced {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

	// Turn Left
	for start, expect := range map[string]string{
		North: West,
		East:  North,
		South: East,
		West:  South,
	} {

		if err := board.Place(robotName, start, 0, 0); err != nil {
			t.Fatalf("Unexpected Error: %s", err.Error())
		}
		if err := board.assetPosition(robotName, start, 0, 0); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := board.Turn(robotName, Left); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := board.assetPosition(robotName, expect, 0, 0); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
	}

	// Turn Right
	for start, expect := range map[string]string{
		North: East,
		East:  South,
		South: West,
		West:  North,
	} {

		if err := board.Place(robotName, start, 0, 0); err != nil {
			t.Fatalf("Unexpected Error: %s", err.Error())
		}
		if err := board.assetPosition(robotName, start, 0, 0); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := board.Turn(robotName, Right); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := board.assetPosition(robotName, expect, 0, 0); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
	}

	if err := board.Turn(robotName, "Madeup"); err == nil {
		t.Errorf("Expected error")
	}
}

func TestExampleCode(t *testing.T) {

	// Testing example C for an 'integration'
	board := NewBoard(4, 4, 1)
	robotName := "foo"

	if err := board.Place(robotName, East, 1, 2); err != nil {
		t.Fatal(err)
	}
	if err := board.Move(robotName); err != nil {
		t.Fatal(err)
	}
	if err := board.Move(robotName); err != nil {
		t.Fatal(err)
	}
	if err := board.Turn(robotName, Left); err != nil {
		t.Fatal(err)
	}
	if err := board.Move(robotName); err != nil {
		t.Fatal(err)
	}

	if err := board.assetPosition(robotName, North, 3, 3); err != nil {
		t.Fatal(err)
	}

}
