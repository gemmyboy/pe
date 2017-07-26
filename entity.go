package pe

import "github.com/gonum/matrix/mat64"

/*
	entity.go
		by Marcus Shannon

	Structure and functionality of an entity.
*/

//Entity Constants
const (
	ShapeCircle    int = 0
	ShapeRectangle int = 1
)

//Entity -: Data structure containing Physics data.
type Entity struct {
	ID int // Unique ID

	//Movement
	LinearVelocity *mat64.Vector // 2D Vector
	Mass           float64
	InverseMass    float64
	LinearPosition *mat64.Vector // 2D Vector
	Force          *mat64.Vector // 2D Vector

	shape        int           // The kind of Rigidbody this is using.
	CenterOfMass *mat64.Vector // Center coordinate of shape - Local Coord

	//Circle
	CircleRadius float64

	//Rectangle
	Min        *mat64.Vector // Bottom left corner - Local Coord
	Max        *mat64.Vector // Top right corner - Local Coord
	rectHeight float64       // internal
	rectWidth  float64       // internal

	//Rotation
	AngularVelocity float64 // Change in current angle
	Inertia         float64
	InverseInertia  float64
	AngularPosition float64 // Current angle
	Torque          float64 // 2D Vector - Rotational Force
} //End Entity

//New -: Create new Entity
func (s *Entity) New() *Entity {
	se := new(Entity)

	//Default Values
	se.LinearVelocity = mat64.NewVector(2, []float64{0, 0})
	se.Mass = 10.0
	se.InverseMass = 1 / se.Mass
	se.LinearPosition = mat64.NewVector(2, []float64{0, 0})
	se.Force = mat64.NewVector(2, []float64{0, 0})

	se.shape = 0
	se.CircleRadius = 10
	se.calcMoment()

	return se
} //End New()

//calcMoment -: called internally to SPEngine, calculates Moment
func (s *Entity) calcMoment() {
	if s.shape == ShapeCircle { //Circle
		s.CenterOfMass = mat64.NewVector(2, []float64{0, 0})
		s.Inertia = (s.Mass * s.CircleRadius * s.CircleRadius) / 4
		s.InverseInertia = 1 / s.Inertia
	} else if s.shape == ShapeRectangle { //Rectangle
		s.rectHeight = s.Min.At(1, 0) - s.Max.At(1, 0)                                    //Calc Height
		s.rectWidth = s.Min.At(0, 1) - s.Max.At(0, 1)                                     //Calc Width
		s.CenterOfMass = mat64.NewVector(2, []float64{s.rectWidth / 2, s.rectHeight / 2}) //Calc COM
		s.Inertia = (s.Mass * (s.rectHeight*s.rectHeight + s.rectWidth*s.rectWidth)) / 12 //Calc Inertia
		s.InverseInertia = 1 / s.Inertia
	}
} //End calcMoment()
