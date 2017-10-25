package robot

import (
	"fmt"
	"testing"
)

func (r Robot) assetPosition(direction string, coords ...int64) error {
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

func TestPlace(t *testing.T) {
	r := Robot{
		Max:       []int64{1, 1},
		Dimension: DirectionSet2D,
	}

	if _, err := r.Report(); err == nil {
		t.Error("Expected error")
	} else if err != ErrorNotPlaced {
		t.Errorf("Wrong error: %s\n", err.Error())
	}

	// Wrong number of coordinates
	if err := r.Place(South); err == nil {
		t.Fatal("Expected error")
	}
	if err := r.Place(South, 0); err == nil {
		t.Fatal("Expected error")
	}
	if err := r.Place(South, 0, 0, 0); err == nil {
		t.Fatal("Expected error")
	}

	// Made up direction
	if err := r.Place("UP", 0, 0); err == nil {
		t.Fatal("Expected error")
	}

	// Fall
	if err := r.Place(South, 2, 0); err == nil {
		t.Fatal("Expected error")
	}
	if err := r.Place(South, 0, 2); err == nil {
		t.Fatal("Expected error")
	}

	if err := r.Place(North, 0, 0); err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}

	reported, err := r.Report()
	if err != nil {
		t.Fatalf("Unexpected Error: %s", err.Error())
	}
	if reported != "0,0,NORTH" {
		t.Errorf("Bad Report: %s", reported)
	}

	if r.Direction.Name != North {
		t.Errorf("Expected North, got %s", r.Direction.Name)
	}

	if len(r.Position) != 2 ||
		r.Position[0] != 0 ||
		r.Position[1] != 0 {
		t.Errorf("Wrong Position: %v", r.Position)
	}
}

func TestMove(t *testing.T) {
	r := Robot{
		Max:       []int64{2, 2},
		Dimension: DirectionSet2D,
	}

	if err := r.Move(); err == nil {
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

		if err := r.Place(direction, 1, 1); err != nil {
			t.Fatalf("Unexpected Error: %s", err.Error())
		}
		if err := r.assetPosition(direction, 1, 1); err != nil {
			t.Fatalf("Testing %s: %s", direction, err.Error())
		}
		if err := r.Move(); err != nil {
			t.Fatalf("Testing %s: %s", direction, err.Error())
		}
		if err := r.assetPosition(direction, position...); err != nil {
			t.Fatalf("Testing %s: %s", direction, err.Error())
		}

		// Try to move once more, would fall
		if err := r.Move(); err == nil {
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

	r := Robot{
		Max:       []int64{0, 0},
		Dimension: DirectionSet2D,
	}

	if err := r.Turn(Left); err == nil {
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

		if err := r.Place(start, 0, 0); err != nil {
			t.Fatalf("Unexpected Error: %s", err.Error())
		}
		if err := r.assetPosition(start, 0, 0); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := r.Turn(Left); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := r.assetPosition(expect, 0, 0); err != nil {
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

		if err := r.Place(start, 0, 0); err != nil {
			t.Fatalf("Unexpected Error: %s", err.Error())
		}
		if err := r.assetPosition(start, 0, 0); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := r.Turn(Right); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
		if err := r.assetPosition(expect, 0, 0); err != nil {
			t.Fatalf("Testing %s: %s", start, err.Error())
		}
	}

	if err := r.Turn("Madeup"); err == nil {
		t.Errorf("Expected error")
	}
}

func TestExampleCode(t *testing.T) {

	// Testing example C for an 'integration'
	r := Robot{
		Max:       []int64{4, 4},
		Dimension: DirectionSet2D,
	}

	if err := r.Place(East, 1, 2); err != nil {
		t.Fatal(err)
	}
	if err := r.Move(); err != nil {
		t.Fatal(err)
	}
	if err := r.Move(); err != nil {
		t.Fatal(err)
	}
	if err := r.Turn(Left); err != nil {
		t.Fatal(err)
	}
	if err := r.Move(); err != nil {
		t.Fatal(err)
	}

	if err := r.assetPosition(North, 3, 3); err != nil {
		t.Fatal(err)
	}

}
