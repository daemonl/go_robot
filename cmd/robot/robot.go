package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/daemonl/go_robot"
)

var maxX int64
var maxY int64
var numRobots int
var quiet bool

func init() {
	flag.Int64Var(&maxX, "max-x", 4, "Maximum board X value")
	flag.Int64Var(&maxY, "max-y", 4, "Maximum board Y value")
	flag.IntVar(&numRobots, "robot-count", 2, "Number of robots on the board")
	flag.BoolVar(&quiet, "quiet", false, "Disable error output (silently ignores bad commands)")
}

func main() {
	flag.Parse()
	board := robot.NewBoard(maxX, maxY, numRobots)

	var errorOutput io.Writer = os.Stdout
	if quiet {
		errorOutput = ioutil.Discard
	}
	if err := robot.CommandStream(board, os.Stdin, os.Stdout, errorOutput); err != nil {
		log.Fatal(err.Error())
	}
}
