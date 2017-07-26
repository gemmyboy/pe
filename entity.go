package pe

import (
	"github.com/gonum/matrix/mat64"
	"github.com/volkerp/goquadtree/quadtree"
)

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
	LinearPosition *mat64.Vector // 2D Vector Also Center of MASS
	Force          *mat64.Vector // 2D Vector

	shape int // The kind of Rigidbody this is using.

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

	//Bounding Box
	BoundingBoxOffSet []float64 //4 indices for xa, xb, ya, yb
} //End Entity

//New -: Create new Entity
func (e *Entity) New(shape int) *Entity {
	se := new(Entity)

	//Default Values
	se.LinearVelocity = mat64.NewVector(2, []float64{0, 0})
	se.Mass = 10.0
	se.InverseMass = 1 / se.Mass
	se.LinearPosition = mat64.NewVector(2, []float64{0, 0})
	se.Force = mat64.NewVector(2, []float64{0, 0})

	se.shape = shape
	se.CircleRadius = 10
	return se
} //End New()

//Calculate -: run entity Calculations
func (e *Entity) Calculate() {
	if e.shape == ShapeCircle { //Circle
		e.LinearPosition = mat64.NewVector(2, []float64{0, 0})
		e.Inertia = (e.Mass * e.CircleRadius * e.CircleRadius) / 4
		e.InverseInertia = 1 / e.Inertia
		e.BoundingBoxOffSet = []float64{e.CircleRadius * -1, e.CircleRadius, e.CircleRadius * -1, e.CircleRadius}

	} else if e.shape == ShapeRectangle { //Rectangle
		e.rectHeight = e.Min.At(1, 0) - e.Max.At(1, 0)                                      //Calc Height
		e.rectWidth = e.Min.At(0, 1) - e.Max.At(0, 1)                                       //Calc Width
		e.LinearPosition = mat64.NewVector(2, []float64{e.rectWidth / 2, e.rectHeight / 2}) //Calc COM
		e.Inertia = (e.Mass * (e.rectHeight*e.rectHeight + e.rectWidth*e.rectWidth)) / 12   //Calc Inertia
		e.InverseInertia = 1 / e.Inertia

		//Bounding Box OffSet
		tc := 0.0
		if e.rectHeight > e.rectWidth {
			tc = e.rectHeight / 2
		} else {
			tc = e.rectWidth / 2
		}
		e.BoundingBoxOffSet = []float64{tc * -1, tc, tc * -1, tc}
	}
} //End Calculate()

//BoundingBox -: implementing method to match type BoundingBoxer in QuadTree
func (e *Entity) BoundingBox() quadtree.BoundingBox {
	return quadtree.NewBoundingBox(
		e.LinearPosition.At(0, 0)+e.BoundingBoxOffSet[0],
		e.LinearPosition.At(0, 0)+e.BoundingBoxOffSet[1],
		e.LinearPosition.At(1, 0)+e.BoundingBoxOffSet[2],
		e.LinearPosition.At(1, 0)+e.BoundingBoxOffSet[3])
} //End BoundingBox()

//SetForce -: Assign Force to the entity
func (e *Entity) SetForce(v *mat64.Vector) {
	e.Force = v
} //End SetForce

//AddForce -: Add Force to the entity
func (e *Entity) AddForce(v *mat64.Vector) {
	e.Force.AddVec(e.Force, v)
} //End AddForce

//SetTorque -: Assign Torque to the entity
func (e *Entity) SetTorque(v float64) {
	e.Torque = v
} //End SetTorque

//AddTorque -: Assign Torque to the entity
func (e *Entity) AddTorque(v float64) {
	e.Torque += v
} //End AddTorque

//ApplyForces -: Apply forces to the entity.
func (e *Entity) ApplyForces(f *mat64.Vector, v float64) {
	if f == nil {
		f = mat64.NewVector(2, nil)
	}
	e.AddForce(f)
	e.AddTorque(v)
} //End ApplyForces

//UpdatePosition -: Update Position of the entity
func (e *Entity) UpdatePosition(dt float64) {

	//Update Velocity | V(t+1) = V(t) + (F(t)*(1/m)*dt)
	tf := mat64.NewVector(2, nil)
	tf.CopyVec(e.Force)
	tf.ScaleVec(e.InverseMass, tf)                // F(t)*(1/m)  <- btw this is acceleration
	tf.ScaleVec(dt, tf)                           // (F(t)*(1/m)*dt) <- splits the acceleration over time
	e.LinearVelocity.AddVec(e.LinearVelocity, tf) // V(t) + (F(t)*(1/m)*dt) <- amount to adjust velocity over time

	//Update Position | P(t+1) = P(t) + (V(t+1)*dt)
	tv := mat64.NewVector(2, nil)
	tv.CopyVec(e.LinearVelocity)
	tv.ScaleVec(dt, tv)                           // (V(t+1)*dt)	<- splits velocity over time
	e.LinearPosition.AddVec(e.LinearPosition, tv) // P(t) + (V(t+1)*dt) <- amount to adjust position over time
} //End UpdatePosition

//UpdateRotation -: Update Rotation of the entity
func (e *Entity) UpdateRotation(dt float64) {

	//Update Angular Velocity | AV(t+1) = AV(t) + (t)*(1/I)*dt
	ta := e.Torque * e.InverseInertia // (t)*(1/I) <- Angular Acceleration
	e.AngularVelocity += (ta * dt)    // AV(t) + (t)*(1/I)*dt <-splits acceleration over time

	//Update Angle | A(t+1) = A(t) + (AV(t+1)*dt)
	e.AngularPosition += (e.AngularVelocity * dt) // A(t) + (AV(t+1)*dt)
} //End UpdateRotation
