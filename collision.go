package pe

import "github.com/volkerp/goquadtree/quadtree"
import "math"

/*
	collision.go
		by Marcus Shannon

	Feature List:
		x AABB vs AABB
		x AABB vs Circle
		x Circle vs Circle
*/

//BroadPair -: list of possible collisions with A
type BroadPair struct {
	a *Entity
	b []quadtree.BoundingBoxer
} //End BroadPair

//NarrowPair -: Pair of colliding Entities
type NarrowPair struct {
	a *Entity
	b *Entity
} //End NarrowPair

//broadphase -: Initial collision detection check
func broadphase(p *PhysicsEngine) []BroadPair {
	n := len(p.entities)
	te := make([]BroadPair, n)

	//Create Quad Tree
	QuadTree := quadtree.NewQuadTree(quadtree.NewBoundingBox(conXA, conXB, conYA, conYB))

	//Quickly Populate the Quad Tree
	for _, e := range p.entities {
		QuadTree.Add(e)
	}

	//Query All to generate possible collisions
	for _, e := range p.entities {
		pairs := QuadTree.Query(e.BoundingBox())
		te = append(te, BroadPair{a: e, b: pairs})
	}

	return te
} //End broadphase()

//narrowphase -: final collision detection check
func narrowphase(p *PhysicsEngine, bp []BroadPair) []NarrowPair {
	n := len(p.entities)
	te := make([]NarrowPair, n*n)
	actual := 0

	for i := 0; i < len(bp); i++ {
		cbp := bp[i]
		a := cbp.a
		for k := 0; k < len(cbp.b); k++ {
			b := cbp.b[k].(*Entity)
			switch cbp.a.shape {
			case ShapeCircle:
				{
					switch b.shape {
					case ShapeCircle:
						{
							if circleVSCircle(a, b) {
								te = append(te, NarrowPair{a: a, b: b})
								actual++
							}
						}
					case ShapeRectangle:
						{
							if circleVSAABB(a, b) {
								te = append(te, NarrowPair{a: a, b: b})
								actual++
							}
						}
					}
				}
			case ShapeRectangle:
				{
					switch b.shape {
					case ShapeCircle:
						{
							if circleVSAABB(a, b) {
								te = append(te, NarrowPair{a: b, b: a})
								actual++
							}
						}
					case ShapeRectangle:
						{
							if aabbVSAABB(a, b) {
								te = append(te, NarrowPair{a: a, b: b})
								actual++
							}
						}
					}
				}
			}
		} //End for
	} //End for
	return te[:actual]
} //End narrowphase()

//circleVSCircle -: Circle vs Circle Collision Check
func circleVSCircle(a *Entity, b *Entity) bool {
	dx := b.LinearPosition.At(0, 0) - a.LinearPosition.At(0, 0)
	dy := b.LinearPosition.At(1, 0) - a.LinearPosition.At(1, 0)
	r := a.CircleRadius + b.CircleRadius

	if (dx*dx)+(dy*dy) > (r * r) {
		return false
	}
	return true
} //End circleVSCircle()

//aabbBSAABB -: Rect vs Rect Collision Check
func aabbVSAABB(a *Entity, b *Entity) bool {
	d1x := b.Min.At(0, 0) - a.Max.At(0, 0)
	d1y := b.Min.At(1, 0) - a.Max.At(1, 0)
	d2x := a.Min.At(0, 0) - b.Max.At(0, 0)
	d2y := a.Min.At(1, 0) - b.Max.At(1, 0)

	if d1x > 0.0 || d1y > 0.0 {
		return false
	} else if d2x > 0.0 || d2y > 0.0 {
		return false
	}
	return true
} //End aabbVSAABB()

//circleVSAABB -: Circle vs AABB Collision Check
func circleVSAABB(a *Entity, b *Entity) bool {
	//Ref: https://stackoverflow.com/questions/21089959/detecting-collision-of-rectangle-with-circle
	distX := math.Abs(a.LinearPosition.At(0, 0) - b.LinearPosition.At(0, 0))
	distY := math.Abs(a.LinearPosition.At(1, 0) - b.LinearPosition.At(1, 0))

	rw := b.rectWidth * float64(0.5)  //Calculation is reused so simplified
	rh := b.rectHeight * float64(0.5) //Calculation is reused so simplified

	if distX > (rw + a.CircleRadius) {
		return false
	}
	if distY > (rh + a.CircleRadius) {
		return false
	}

	if distX <= rw {
		return true
	}
	if distY <= rh {
		return true
	}

	dx := distX - rw
	dy := distY - rh
	return (dx*dx + dy*dy) <= (a.CircleRadius * a.CircleRadius)
} //End circleVSAABB()
