package models

type WorkoutDetail struct {
	Id             uint
	SessionId      int
	MachineId      int
	MinuteDuration int
	Sets           int
	Reps           int
}
