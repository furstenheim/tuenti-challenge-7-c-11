package main

import "fmt"
import "./priority-queue"

type galaxy struct {
	id uint16
	distances map[uint]uint32
	charges map[uint]uint32
	visited map[uint]bool
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
	addChargeToGalaxy(&mUniverse, 0, "Green", 15)

	fmt.Println(mUniverse.galaxies[0])
	fmt.Println(mUniverse)
	v1 := visit{color: 3, distance: 1, justArrived: false, galaxy: 1}
	v2 := visit{color: 1, distance: 4, justArrived: false, galaxy: 1}
	mPq := pq.New()
	insertVisit(mPq, v1)
	insertVisit(mPq, v2)
	v3,_ := mPq.Pop()
	fmt.Println(v3)

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
	mGalaxy := galaxy{id: id, distances: make(map[uint]uint32), charges: make(map[uint]uint32), visited: make(map[uint]bool)}
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
func insertVisit (pq pq.PriorityQueue, visit visit) {
	pq.Insert(visit, visit.distance) // probably should consider color
}

