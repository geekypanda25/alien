package alienworld

import (
	
	"alien/queue"

	"errors"
	"fmt"
	"log"
	"strings"

)

const (
	// the maximum number of aliens at a city.
	maxaliens = 2
	// maxpaths is 4 ie NSEW
	maxpaths = 4
)

//map struct to hold cities and aliens
type Map struct {

	cities map[string]*City

	aliens map[string]*Alien

}

//City struct holds name, paths in and out, and alien in the city
type City struct {
	
	name     string

	inpaths  map[string]string
	
	outpaths map[string]string
	
	alienAt  map[string]*Alien
}

// Priority implements the Heapable interface.
func (c *City) Priority(other interface{}) bool {
	
	if t, ok := other.(*City); ok {
	
		return len(c.outpaths) > len(t.outpaths)
	
	}

	return false
}

// String interface
func (c *City) String() string {
	
	if len(c.outpaths) == 0 {
	
		return ""
	
	}

	paths := ""

	for pathdir, pathtocity := range c.outpaths {
	
		paths += fmt.Sprintf(" %s=%s", pathdir, pathtocity)
	
	}

	return fmt.Sprintf("%s%s", c.name, paths)
}

// Returns a new Map.
func NewMap() *Map {
	
	return &Map{
	
		cities: make(map[string]*City),
	
		aliens: make(map[string]*Alien),
	
	}
}

// aNames returns a list of all aliens on map
func (m *Map) ANames() []string {
	
	aNames := make([]string, 0, len(m.aliens))

	for aName := range m.aliens {
	
		aNames = append(aNames, aName)
	
	}

	return aNames
}

// Cities on the map
func (m *Map) Cities() []*City {
	
	cities := make([]*City, 0, len(m.cities))

	for _, city := range m.cities {
	
		cities = append(cities, city)
	
	}

	return cities
}

// Number of Cities on the map.
func (m *Map) NumCities() uint {
	
	return uint(len(m.cities))

}


// Aliens on the map
func (m *Map) NumAliens() uint {

	return uint(len(m.aliens))

}

// Returns all cities on map
func (m *Map) cNames() []string {

	cNames := make([]string, 0, m.NumCities())

	for cName := range m.cities {

		cNames = append(cNames, cName)

	}

	return cNames
}

// addpaths adds a path(directional edge) from an origin city to a linked city.
func (m *Map) Addpaths(cName, linkCityDir, pathtocity string) {
	// Add origin city
	if _, ok := m.cities[cName]; !ok {

		m.cities[cName] = &City{
			name:     cName,

			inpaths:  make(map[string]string, maxpaths),
			
			alienAt:  make(map[string]*Alien, maxaliens),

			outpaths: make(map[string]string, maxpaths),
		}
	}

	// Add linked city
	if _, ok := m.cities[pathtocity]; !ok {

		m.cities[pathtocity] = &City{
			name:     pathtocity,

			inpaths:  make(map[string]string, maxpaths),
			
			alienAt:  make(map[string]*Alien, maxaliens),

			outpaths: make(map[string]string, maxpaths),
		}
	}

	// Added in and out paths
	m.cities[pathtocity].inpaths[strings.ToLower(linkCityDir)] = cName

	m.cities[cName].outpaths[strings.ToLower(linkCityDir)] = pathtocity
	
}

// Moves alien on the map and fights them if there are more than maxAliens in a city
func (m *Map) MoveAlien() (string, error) {

	for _, alien := range m.aliens {

		occupiedCity := alien.cName

		city := m.cities[occupiedCity]

		for _, pathtocity := range city.outpaths {

			linkCity := m.cities[pathtocity]

			if len(linkCity.alienAt) < maxaliens {

				delete(city.alienAt, alien.name)

				alien.cName = linkCity.name

				linkCity.alienAt[alien.name] = alien

				return alien.name, nil
			}
		}
	}

	return "", errors.New("unable to move any alien")
}

// removes a destroyed city
func (m *Map) destroyCity(city *City) []string {

	destroyedAliens := make([]string, 0, maxaliens)

	for aName := range city.alienAt {

		destroyedAliens = append(destroyedAliens, aName)

		delete(m.aliens, aName)

	}

	for _, inlinkCity := range city.inpaths {

		inCity := m.cities[inlinkCity]

		for pathdir, pathtocity := range inCity.outpaths {

			if pathtocity == city.name {

				delete(inCity.outpaths, pathdir)

				break
			}
		}

		for pathdir, pathtocity := range inCity.inpaths {

			if pathtocity == city.name {

				delete(inCity.inpaths, pathdir)

				break
			}
		}
	}

	for _, outlinkCity := range city.outpaths {

		outCity := m.cities[outlinkCity]

		for pathdir, pathtocity := range outCity.inpaths {

			if pathtocity == city.name {

				delete(outCity.inpaths, pathdir)

				break
			}
		}
	}

	delete(m.cities, city.name)

	return destroyedAliens
}

// If 2 aliens are at city it fights them, in the process destroying the city
func (m *Map) ExecuteFights() {

	for _, alien := range m.aliens {

		occupiedCity := alien.cName

		city := m.cities[occupiedCity]

		if len(city.alienAt) == maxaliens {

			destroyedAliens := m.destroyCity(city)

			log.Printf("%s has been destroyed by %s!", city.name, strings.Join(destroyedAliens, " and "))
		}
	}
}

// Adds n aliens to map
func (m *Map) AddAliens(n uint) {

	pq := queue.NewPriorityQueue()

	//adds cities based on number of out links, highest number of paths are prioritized
	for _, city := range m.cities {

		pq.Push(city)
	}

	seededAliens := uint(0)

	for seededAliens != n {

		city := pq.Pop().(*City)

		for i := 0; i < maxaliens && seededAliens != n; i++ {

			alien := &Alien{

				name:     fmt.Sprintf("alien%d", seededAliens+1),

				cName: city.name,
			}

			city.alienAt[alien.name] = alien

			m.aliens[alien.name] = alien

			seededAliens++
		}
	}
}

// Stringer interface.
func (m *Map) String() (s string) {

	for _, city := range m.cities {

		aliens := make([]string, 0, len(city.alienAt))

		for alien := range city.alienAt {

			aliens = append(aliens, alien)
		}

		s += fmt.Sprintf("{city: %s, outpaths: %s, inpaths: %s, alienAt: [%s]}\n",city.name, city.outpaths, city.inpaths, strings.Join(aliens, " "),

		)
	}

	return
}
