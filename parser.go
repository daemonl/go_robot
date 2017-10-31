package robot

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type iBoard interface {
	Place(robot string, direction string, position ...int64) error
	Move(robot string) error
	Turn(robot string, direction string) error
	Report(robot string) (string, error)
}

// DoCommand parses the given command and applies it to the robot, returning
// any errors or command output
func DoCommand(robot iBoard, rawCommand string) (string, error) {
	if robot == nil {
		return "", fmt.Errorf("No robot configured")
	}

	command, err := parseCommand(rawCommand)
	if err != nil {
		return "", err
	}

	switch command.command {
	case "PLACE":
		if len(command.args) < 2 {
			return "", fmt.Errorf("PLACE requires at least 2 arguments")
		}
		direction := command.args[len(command.args)-1]
		coordStrings := command.args[0 : len(command.args)-1]
		coords := make([]int64, len(coordStrings), len(coordStrings))
		for idx, raw := range coordStrings {
			out, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return "", fmt.Errorf("parsing parameter %d: %s", idx, err.Error())
			}
			coords[idx] = out
		}
		return "", robot.Place(command.robotName, direction, coords...)

	case "MOVE":
		if len(command.args) != 0 {
			return "", fmt.Errorf("MOVE does not take any arguments")
		}
		return "", robot.Move(command.robotName)

	case "REPORT":
		if len(command.args) != 0 {
			return "", fmt.Errorf("MOVE does not take any arguments")
		}
		return robot.Report(command.robotName)

	case "LEFT", "RIGHT":
		if len(command.args) != 0 {
			return "", fmt.Errorf("%s does not take any arguments", command.command)
		}
		return "", robot.Turn(command.robotName, command.command)

	default:
		return "", fmt.Errorf("No such command '%s'", command.command)
	}

}

// Split at any combo of space and comma, but only one comma
var reSplit = regexp.MustCompile(`([ ]*,[ ]*|[ ]+)`)

type parsedCommand struct {
	robotName string
	command   string
	args      []string
}

func parseCommand(raw string) (*parsedCommand, error) {
	parts := reSplit.Split(raw, -1)
	if len(parts) == 1 {
		return nil, fmt.Errorf("Invalid command, robot name required")

	}

	if len(parts) == 2 {
		return &parsedCommand{
			robotName: parts[1],
			command:   parts[0],
		}, nil
	}

	return &parsedCommand{
		robotName: parts[len(parts)-1],
		command:   parts[0],
		args:      parts[1 : len(parts)-1],
	}, nil
}

// CommandStream issues a series of commands to the robot, displaying errors if showErrors is true
func CommandStream(robot iBoard, streamIn io.Reader, streamOut io.Writer, streamError io.Writer) error {
	scanner := bufio.NewScanner(streamIn)
	for scanner.Scan() {
		command := scanner.Text()
		output, err := DoCommand(robot, command)
		if err != nil {
			streamError.Write([]byte(err.Error() + "\n"))
		}
		if output != "" {
			streamOut.Write([]byte(output + "\n"))
		}
	}
	return scanner.Err()
}
