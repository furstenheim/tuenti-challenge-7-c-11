package main

import "fmt"
import "./priority-queue"

type galaxy struct {
	id uint16
	distances map[uint]uint32
	charges map[uint]uint32
	visited map[uint]bool
	wormholes map[uint]([]uint) // galaxy and number of possible colors
}

type visit struct {
	color uint16
	distance uint32
	justArrived bool
	galaxy uint16
}

type universe struct {
	galaxies []galaxy
	visits pq.PriorityQueue
	color2id map[string]uint
	primary2shift map[string]uint
	numPrimCol uint
}

func main() {
	mUniverse := universe{galaxies: []galaxy{}, visits: pq.New(), color2id: make(map[string]uint), primary2shift: make(map[string]uint), numPrimCol: 0}
	addColorToUniverse(&mUniverse, "Red", 0, []string{} )
	addColorToUniverse(&mUniverse, "Blue", 0, []string{} )
	addColorToUniverse(&mUniverse, "Green", 0, []string{} )
	addColorToUniverse(&mUniverse, "Yellow", 2, []string{"Red", "Green"} )
	createGalaxy(&mUniverse)
	addChargeToGalaxy(&mUniverse, 0, "Red", 10)
	addChargeToGalaxy(&mUniverse, 0, "Green", 10)
	createGalaxy(&mUniverse)
	addChargeToGalaxy(&mUniverse, 1, "Red", 15)
	addChargeToGalaxy(&mUniverse, 1, "Green", 15)
	createGalaxy(&mUniverse)
	addChargeToGalaxy(&mUniverse, 2, "Blue", 7)
	addChargeToGalaxy(&mUniverse, 2, "Green", 5)
	createGalaxy(&mUniverse)
	addChargeToGalaxy(&mUniverse, 3, "Red", 11)
	createGalaxy(&mUniverse)
	addChargeToGalaxy(&mUniverse, 4, "Red", 10)
	addChargeToGalaxy(&mUniverse, 4, "Green", 10)
	addWormHole(&mUniverse, "Red", 0, 1)
	addWormHole(&mUniverse, "Red", 3, 2)
	addWormHole(&mUniverse, "Green", 1, 2)
	addWormHole(&mUniverse, "Blue", 0, 3)
	addWormHole(&mUniverse, "Blue", 3, 4)
	addWormHole(&mUniverse, "Yellow", 2, 0)
	v1 := visit{color: 0, distance: 0, justArrived: true, galaxy: 0}

	insertVisit(mUniverse.visits, v1)

	v3,err := mUniverse.visits.Pop()
	fmt.Println(v3, err)
}

func addColorToUniverse (universe *universe, name string, number uint, primaries []string) {
	var newid uint
	if (number == 0) {
		newid = 1 << universe.numPrimCol
		universe.primary2shift[name] = universe.numPrimCol
		fmt.Println(universe.numPrimCol)
		universe.numPrimCol = universe.numPrimCol + 1
		fmt.Println(universe.numPrimCol)
		universe.color2id[name] = newid
	} else {
		newid = 0
		for _,primary := range primaries {
			shift := universe.primary2shift[primary]
			newid += (1 << shift)
		}
		universe.color2id[name] = newid
	}
}

func createGalaxy (universe *universe) {
	id := uint16(len(universe.galaxies))
	mGalaxy := galaxy{id: id, distances: make(map[uint]uint32), charges: make(map[uint]uint32), visited: make(map[uint]bool), wormholes: make(map[uint]([]uint))}
	for _, color := range universe.color2id {
		mGalaxy.visited[color] = false
	}
	universe.galaxies = append(universe.galaxies, mGalaxy)
}
func addChargeToGalaxy (universe *universe, galaxyId uint, color string, time uint32) {
	id := universe.color2id[color]
	galaxy := universe.galaxies[galaxyId]
	galaxy.charges[id] = time
}

func addWormHole (universe *universe, color string, start uint, end uint) {
	id := universe.color2id[color]
	galaxy := universe.galaxies[start]

	_, ok := galaxy.wormholes[end]
	if (ok) {
		galaxy.wormholes[end] = []uint{id}
	} else {
		galaxy.wormholes[end] = append(galaxy.wormholes[end], id)
	}
}
func insertVisit (pq pq.PriorityQueue, visit visit) {
	pq.Insert(visit, visit.distance) // probably should consider color
}

