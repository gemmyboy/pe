package pe

/*
	physicsengine.go
		by Marcus Shannon

	Package is geared towards processing physics based information.  It's meant to be integrated with TSP for networking purposes.append
	Ideally, this package is optimized for quick networking capabilities.

	Features:
	- Entity Update Processing
		- Input Processing
		- Entity Add/Delete
	- Collision Detection
		- Detect Circle - Circle
		- Detect Square - Square
		- Detect Circle - Square
	- Network Events
		- Standardize Event Processing Flagging
		- Slim down structures for portability
		- Create Group in SS with Game Engine
*/
