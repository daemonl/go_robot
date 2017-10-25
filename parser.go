package robot

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type IRobot interface {
	Place(string, ...int64) error
	Move() error
	Turn(string) error
	Report() (string, error)
}

func DoCommand(robot IRobot, command string) (string, error) {
	if robot == nil {
		return "", fmt.Errorf("No robot configured")
	}

	command, args := parseCommand(command)

	switch command {
	case "PLACE":
		if len(args) < 2 {
			return "", fmt.Errorf("PLACE requires at least 2 arguments")
		}
		direction := args[len(args)-1]
		coordStrings := args[0 : len(args)-1]
		coords := make([]int64, len(coordStrings), len(coordStrings))
		for idx, raw := range coordStrings {
			out, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return "", fmt.Errorf("parsing parameter %d: %s", idx, err.Error())
			}
			coords[idx] = out
		}
		return "", robot.Place(direction, coords...)

	case "MOVE":
		if len(args) != 0 {
			return "", fmt.Errorf("MOVE does not take any arguments")
		}
		return "", robot.Move()

	case "REPORT":
		if len(args) != 0 {
			return "", fmt.Errorf("MOVE does not take any arguments")
		}
		return robot.Report()

	case "LEFT", "RIGHT":
		if len(args) != 0 {
			return "", fmt.Errorf("%s does not take any arguments", command)
		}
		return "", robot.Turn(command)

	default:
		return "", fmt.Errorf("No such command '%s'", command)
	}

}

// Split at any combo of space and comma, but only one comma
var reSplit = regexp.MustCompile(`([ ]*,[ ]*|[ ]+)`)

func parseCommand(raw string) (cmd string, args []string) {
	parts := reSplit.Split(raw, -1)
	if len(parts) == 1 {
		return parts[0], []string{}
	}

	return parts[0], parts[1:]
}

// CommandStream issues a series of commands to the robot, displaying errors if showErrors is true
func CommandStream(robot IRobot, streamIn io.Reader, streamOut io.Writer, streamError io.Writer) error {
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
