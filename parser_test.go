package robot

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {

	for raw, expect := range map[string]parsedCommand{
		"MOVE FOO": {
			command: "MOVE",
			args:    []string{},
		},
		"PLACE 1,2,LEFT FOO": {
			command: "PLACE",
			args:    []string{"1", "2", "LEFT"},
		},
		"PLACE 1,2    ,LEFT, FOO": {
			command: "PLACE",
			args:    []string{"1", "2", "LEFT"},
		},
	} {
		out, err := parseCommand(raw)
		if err != nil {
			t.Errorf("Bad parse: %s", err.Error())
			continue
		}
		if out.command != expect.command {
			t.Errorf("Bad parse %s -> '%s'", raw, out)
			continue
		}
		if len(expect.args) != len(out.args) {
			t.Errorf("Bad parse %s -> %v", raw, out.args)
		} else {
			for idx := range expect.args {
				if expect.args[idx] != out.args[idx] {
					t.Errorf("Bad parse %s -> %v", raw, out.args)
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

func (cr *captureRobot) Place(robot string, direction string, pos ...int64) error {
	cr.logf("PLACE %v %s", pos, direction)
	return nil
}

func (cr *captureRobot) Move(robot string) error {
	cr.logf("MOVE")
	return nil
}

func (cr *captureRobot) Turn(robot string, direction string) error {
	cr.logf("%s", direction)
	return nil
}

func (cr *captureRobot) Report(robot string) (string, error) {
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
		"PLACE 1,1,NORTH FOO",
		"PLACE 1,NORTH FOO",
		"PLACE 1,1,1,NORTH FOO",
		"MOVE FOO",
		"REPORT FOO",
		"LEFT FOO",
	} {
		if _, err := DoCommand(capture, cmd); err != nil {
			t.Fatal(err.Error())
		}
	}

	// BAD
	for _, cmd := range []string{
		"PLACE FOO",
		"PLACE A,1,NORTH FOO",
		"MOVE 1 FOO",
		"REPORT SOMETHING FOO",
		"FOOBAR FOO",
		"LEFT 1 FOO",
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

	board := NewBoard(5, 5, 1)

	commands := strings.Join([]string{
		"PLACE 1,2,EAST FOO",
		"MOVE FOO",
		"MOVE FOO",
		"LEFT FOO",
		"MOVE FOO",
		"REPORT FOO",
		"ERROR FOO",
	}, "\n")

	// Test, not including errors in output
	bufferOut := bytes.NewBuffer([]byte{})
	bufferError := bytes.NewBuffer([]byte{})
	bufferIn := bytes.NewBufferString(commands)
	if err := CommandStream(board, bufferIn, bufferOut, bufferError); err != nil {
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
