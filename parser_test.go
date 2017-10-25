package robot

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {

	for raw, expect := range map[string]struct {
		command string
		args    []string
	}{
		"MOVE": {
			command: "MOVE",
			args:    []string{},
		},
		"PLACE 1,2,LEFT": {
			command: "PLACE",
			args:    []string{"1", "2", "LEFT"},
		},
		"PLACE 1,2    ,LEFT": {
			command: "PLACE",
			args:    []string{"1", "2", "LEFT"},
		},
	} {
		out, args := parseCommand(raw)
		if out != expect.command {
			t.Errorf("Bad parse %s -> '%s'", raw, out)
			continue
		}
		if len(expect.args) != len(args) {
			t.Errorf("Bad parse %s -> %v", raw, args)
		} else {
			for idx := range expect.args {
				if expect.args[idx] != args[idx] {
					t.Errorf("Bad parse %s -> %v", raw, args)
				}
			}
		}
	}
}

type captureRobot struct {
	commands []string
}

func (cr *captureRobot) logf(f string, args ...interface{}) {
	cr.commands = append(cr.commands, fmt.Sprintf(f, args...))
}

func (cr *captureRobot) Place(direction string, pos ...int64) error {
	cr.logf("PLACE %v %s", pos, direction)
	return nil
}

func (cr *captureRobot) Move() error {
	cr.logf("MOVE")
	return nil
}

func (cr *captureRobot) Turn(direction string) error {
	cr.logf("%s", direction)
	return nil
}

func (cr *captureRobot) Report() (string, error) {
	cr.logf("REPORT")
	return "FAKE", nil
}

func (cr *captureRobot) assert(expect []string) error {
	if len(expect) != len(cr.commands) {
		return fmt.Errorf("Wrong Length %#v", cr.commands)
	}
	for idx, expectVal := range expect {
		if expectVal != cr.commands[idx] {
			return fmt.Errorf("at index %d: %s", idx, cr.commands[idx])
		}
	}
	return nil
}

func TestCommander(t *testing.T) {

	capture := &captureRobot{
		commands: []string{},
	}

	if _, err := DoCommand(nil, "REPORT"); err == nil {
		t.Fatal("Expected error")
	}

	for _, cmd := range []string{
		"PLACE 1,1,NORTH",
		"PLACE 1,NORTH",
		"PLACE 1,1,1,NORTH",
		"MOVE",
		"REPORT",
		"LEFT",
	} {
		if _, err := DoCommand(capture, cmd); err != nil {
			t.Fatal(err.Error())
		}
	}

	// BAD
	for _, cmd := range []string{
		"PLACE",
		"PLACE A,1,NORTH",
		"MOVE 1",
		"REPORT SOMETHING",
		"FOOBAR",
		"LEFT 1",
	} {
		if _, err := DoCommand(capture, cmd); err == nil {
			t.Fatal("Expected Error")
		}
	}
	if err := capture.assert([]string{
		"PLACE [1 1] NORTH",
		"PLACE [1] NORTH",
		"PLACE [1 1 1] NORTH",
		"MOVE",
		"REPORT",
		"LEFT",
	}); err != nil {
		t.Fatal(err.Error())
	}
}

func TestStream(t *testing.T) {

	r := &Robot{
		Dimension: DirectionSet2D,
		Max:       []int64{5, 5},
	}

	commands := strings.Join([]string{
		"PLACE 1,2,EAST",
		"MOVE",
		"MOVE",
		"LEFT",
		"MOVE",
		"REPORT",
		"ERROR",
	}, "\n")

	// Test, not including errors in output
	bufferOut := bytes.NewBuffer([]byte{})
	bufferError := bytes.NewBuffer([]byte{})
	bufferIn := bytes.NewBufferString(commands)
	if err := CommandStream(r, bufferIn, bufferOut, bufferError); err != nil {
		t.Fatal(err.Error())
	}

	output := bufferOut.String()
	if output != "3,3,NORTH\n" {
		t.Error(output)
	}

	output = bufferError.String()
	if output != "No such command 'ERROR'\n" {
		t.Error(output)
	}

}
