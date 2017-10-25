Toy Robot Simulator
===================

This is a response to a coding test. See SPEC.md for more details on the test.

The application is a simulation of a toy robot moving on a square tabletop.

Running
-------

The application can be built using Go version 1.3 or above. There are no
external dependencies.

Running directly with go allows for flags.
```
go run ./cmd/robot/*.go
```

Some makefile shortcuts are available:
```
make run
```

To run a spec file from the example folder:
```
make runspec/spec-a
```
or
```
cat ./example/spec-a.txt | make run
```


