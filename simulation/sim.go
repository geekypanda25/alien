package simulation

import (
	"alien/alienworld"
)

const (

	minMoves = 1000

)

// Simulation reflects a simulation of an alien invasion on a given world map.
type Simulation struct {

	alienMap   *alienworld.Map

	alienMoves map[string]uint

}

// Returns a newly initialized sim
func Newsim(alienMap *alienworld.Map) *Simulation {

	s := &Simulation{

		alienMap:   alienMap,

		alienMoves: make(map[string]uint),

	}

	for _, alienName := range alienMap.ANames() {

		s.alienMoves[alienName] = 0

	}

	return s
}

// Executes a simulation.
// 1) Execute random moves and Run Fights.
// 2) Tracks number of moves for alien.
// 3) Sim ends when either a) All aliens are destroyed
//						   b) minMoves is reached
// An error is returned if the simulation fails.
func (s *Simulation) Run() error {

	for s.canmove() {

		alienName, err := s.alienMap.MoveAlien()

		if err != nil {

			return err

		}

		_, ok := s.alienMoves[alienName]

		if ok {

			s.alienMoves[alienName]++

			if s.alienMoves[alienName] >= minMoves {

				delete(s.alienMoves, alienName)

			}
		}

		s.alienMap.ExecuteFights()
	}

	return nil
}

// Return a boolean on whether or not a simulation can continue
func (s *Simulation) canmove() bool {

	if s.alienMap.NumAliens() == 0 {

		return false

	}

	for _, totalMoves := range s.alienMoves {

		if totalMoves < minMoves {

			return true

		}
	}

	return false
}
