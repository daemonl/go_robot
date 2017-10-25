package main

import (
	"log"
	"os"

	robot "github.com/daemonl/go_robot"
)

func main() {
	r := &robot.Robot{
		Dimension: robot.DirectionSet2D,
		Max:       []int64{5, 5},
	}

	if err := robot.CommandStream(r, os.Stdin, os.Stdout, os.Stderr); err != nil {
		log.Fatal(err.Error())
	}
}
