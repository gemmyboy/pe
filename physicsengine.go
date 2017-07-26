package pe

import (
	"log"
	"sync"
	"time"

	"github.com/gonum/matrix/mat64"
)

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
		x Detect Circle - Circle
		x Detect Square - Square
		x Detect Circle - Square
	- Network Events
		- Standardize Event Processing Flagging
		x Slim down structures for portability
		- Create Group in SS with Game Engine
*/

const (
	conEventChannelLimit  = 10000 //Queue Memory Limit
	conForcesChannelLimit = 10000 //Queue Memory Limit

	conEntityLimit = 100000 //Initial Memory Map Allocation

	//World Bounding Box Coordinates for BroadPhaseQuadTree -- Default
	conXA = float64(-100000)
	conXB = float64(100000)
	conYA = float64(-100000)
	conYB = float64(100000)
)

//PhysicsEngine -: data structure
type PhysicsEngine struct {

	//Entities
	entities map[int]*Entity

	//Time
	timePreviousStep time.Time
	timeCurrentStep  time.Time
	timeDelta        time.Duration
	timeDT           float64

	//Event Management
	eventChannel chan *Note  //Queue of Notes to process
	eventMux     *sync.Mutex //Mutex to cut off Channel

	//Apply Forces
	forcesChannel chan *Note //Queue of Notes to process

} //End PhysicsEngine

//New -: Create a new PhysicsEngine
func New() *PhysicsEngine {
	p := new(PhysicsEngine)

	//Entities
	p.entities = make(map[int]*Entity, conEntityLimit)

	//Time Step
	p.timePreviousStep = time.Now()

	//Event Management
	p.eventChannel = make(chan *Note, conEventChannelLimit)
	p.eventMux = &sync.Mutex{}

	//Apply Forces
	p.forcesChannel = make(chan *Note, conForcesChannelLimit)

	return p
} // End New()

//Step -: Take a single step in time for simulation
func (p *PhysicsEngine) Step() {

	p.PreTime()

	p.eventManagement()
	p.applyForces()
	p.updateEntities()
	np := p.collisionHandling()
	p.solveConstraints(np)

	p.PostTime()
} //End Step()

//---------------------------------------------------------------------

//PreTime -: Run time calculations
func (p *PhysicsEngine) PreTime() {
	p.timeCurrentStep = time.Now()
	p.timeDelta = p.timeCurrentStep.Sub(p.timePreviousStep)
	p.timeDT = p.timeDelta.Seconds()
} //End PreTime()

//PostTime -: Finish time calculations
func (p *PhysicsEngine) PostTime() {
	p.timePreviousStep = p.timeCurrentStep
} //End PostTime()

//eventManagement -: Process incoming events
func (p *PhysicsEngine) eventManagement() {
	p.eventMux.Lock()
	for {
		select {
		case n := <-p.eventChannel: //Non-blocking Channel
			{
				if n.From == IDGameEngine {
					p.eventGameEngine(n)
				} else if n.From == IDClient {
					p.eventClient(n)
				} else {
					continue
				}
			}
		default:
			{
				p.eventMux.Unlock()
				return
			}
		}
	} //End for
} //End eventManagement()

//appleForces -: Update force & torque in all changed entities
func (p *PhysicsEngine) applyForces() {
	for {
		select {
		case n := <-p.forcesChannel: //Non-blocking Channel
			{
				if n.Flag == FlagUpdateForces {
					id := n.Data[0].(int)
					f := n.Data[1].(*mat64.Vector)
					t := n.Data[2].(float64)
					if e, ok := p.entities[id]; ok {
						e.ApplyForces(f, t)
					} else {
						log.Println("WARNING: Bad ID; Entity ID:", id, "does not exist!")
					}
				} else {
					log.Println("WARNING: Bad Flag; ApplyForces Flag:", n.Flag)
					continue
				}
			}
		default:
			{
				return
			}
		}
	} //End for
} //End applyForces

//updateEntities -: Update entities based on forces
func (p *PhysicsEngine) updateEntities() {
	for _, e := range p.entities {
		e.UpdatePosition(p.timeDT)
		e.UpdateRotation(p.timeDT)
	} //End for
} //End updateEntities()

//collisionHandling -: Generate Collision Pairs
func (p *PhysicsEngine) collisionHandling() []NarrowPair {
	bp := broadphase(p)
	np := narrowphase(p, bp)
	return np
} //End collisionHandling()

//solveConstraints -:
func (p *PhysicsEngine) solveConstraints(np []NarrowPair) {

	//Physics basically says that: For Elastic Collisions, just exchange the Velocity
	//	vector between the 2 entities Colliding.
	//	Note: We can add friction/dampening to slow it down.

	for i := 0; i < len(np); i++ {
		tv := np[i].a.LinearVelocity
		np[i].a.LinearVelocity = np[i].b.LinearVelocity
		np[i].b.LinearVelocity = tv
	} //End for

	//This is probably more complicated but it should work for now.
} //End solveConstraints()

//-------------------------------------------------------------------------------

//eventGameEngine -: Process Note from the GameEngine
func (p *PhysicsEngine) eventGameEngine(n *Note) {
	switch n.Flag {
	case FlagSpawn:
		{

		}
	case FlagKill:
		{

		}
	default:
		{
			log.Println("WARNING: Bad Note; GameEngine Flag:", n.Flag)
		}
	} //End Switch
} //End eventGameEngine()

//eventClient -: Process Note from the Client
func (p *PhysicsEngine) eventClient(n *Note) {
	switch n.Flag {
	case FlagInput:
		{

		}
	default:
		{
			log.Println("WARNING: Bad Note; GameEngine Flag:", n.Flag)
		}
	} //End Switch
} //End eventClient()
