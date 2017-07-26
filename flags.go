package pe

/*
	flags.go
		by Marcus Shannon

	List of constants used through-out to signal flagging with Note(s).
*/

//Flags for what the note is for
const (
	FlagCollision uint32 = 0
	FlagInput     uint32 = 1
	FlagSpawn     uint32 = 2
	FlagKill      uint32 = 3

	/*
		FlagUpdateForces - Update the Forces in PhysicsEngine
		data[0] = int				//Entity ID
		data[1] = *mat64.Vector		//Force Vector
		data[2] = float64			//Torque
	*/
	FlagUpdateForces uint32 = 4

	//Identifiers for each System
	IDGameEngine    uint32 = 4862
	IDPhysicsEngine uint32 = 5913
	IDClient        uint32 = 8177
)
