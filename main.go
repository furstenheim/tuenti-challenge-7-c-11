package main

import "fmt"
import "./priority-queue"


func main() {
	mUniverse := universe{galaxies: []galaxy{}, visits: pq.New(), color2id: make(map[string]uint), primary2shift: make(map[string]uint), numPrimCol: 0}
	mUniverse.color2id["Void"] = 0
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
	wrapUpGalaxy(&mUniverse, 0)
	wrapUpGalaxy(&mUniverse, 1)
	wrapUpGalaxy(&mUniverse, 2)
	wrapUpGalaxy(&mUniverse, 3)
	wrapUpGalaxy(&mUniverse, 4)

	mUniverse = travelUniverse(mUniverse)
	fmt.Println(formatUniverse(mUniverse))
}

func travelUniverse (universe universe) universe {
	fmt.Println(universe.color2id)
	v0 := visit{color: 0, distance: 0, justArrived: true, galaxy: 0}
	insertVisit(universe.visits, v0)
	v1, areThereVisits := universe.visits.Pop()
	for areThereVisits == nil {
		newVisit, _ := v1.(visit)
		// fmt.Println(newVisit)
		wasUseful := visitTheGalaxy(&universe, newVisit)
		if (wasUseful) {
			nextVisits := planNewVisits(universe, newVisit)
			for _, nextVisit := range(nextVisits) {
				insertVisit(universe.visits, nextVisit)
			}
		}

		v1, areThereVisits = universe.visits.Pop()
	}
	return universe
}

func visitTheGalaxy (universe *universe, visit visit) bool {
	galaxy := universe.galaxies[visit.galaxy]
	distance, isVisited := galaxy.distances[visit.color]
	if (!isVisited || (distance > visit.distance)) {
		// Initiate all subcolors
		for _, color := range universe.color2id {
			if (contained(color, visit.color)) {
				colorDistance, wasVisited := galaxy.distances[color]
				if (!wasVisited || (colorDistance > visit.distance)) {
					galaxy.distances[color] = visit.distance
				}
			}
		}
		return true
	}
	return true
}

func planNewVisits (universe universe, visit visit) []visit {
	galaxy := universe.galaxies[visit.galaxy]
	stays := getCharges(galaxy, visit)
	travels := getWormHoles(universe, galaxy, visit)
	return append(stays, travels...)
}

func getCharges (galaxy galaxy, originalVisit visit) []visit {
	originalColor := originalVisit.color
	stays := []visit{}
	/*if (!originalVisit.justArrived) {
		return stays
	}*/
	for color, distance := range(galaxy.charges) {
		if (!contained(color, originalColor)) {
			newVisit := visit{color: color | originalColor, distance: originalVisit.distance + distance, justArrived: false, galaxy: galaxy.id}
			galaxyDistance, isVisited := galaxy.distances[newVisit.color]
			if (!isVisited || (galaxyDistance > newVisit.distance)) {
				stays = append(stays, newVisit)
			}
		}
	}
	return stays
}

func getWormHoles (universe universe, galaxy galaxy, originalVisit visit) []visit {
	travels := []visit{}

	for destinationId, wormholes := range galaxy.wormholes {
		destinationGalaxy := universe.galaxies[destinationId]
		for color, _ := range wormholes {
			if (contained(color, originalVisit.color)) {
				newVisit := visit{color: originalVisit.color ^ color, distance: originalVisit.distance, justArrived: true, galaxy: destinationId}
				galaxyDistance, isVisited := destinationGalaxy.distances[newVisit.color]
				if (!isVisited || (galaxyDistance > newVisit.distance)) {
					travels = append(travels, newVisit)
				}

			}
		}
	}
	return travels
}

func addColorToUniverse (universe *universe, name string, number uint, primaries []string) {
	var newid uint
	if (number == 0) {
		newid = 1 << universe.numPrimCol
		universe.primary2shift[name] = universe.numPrimCol
		universe.numPrimCol = universe.numPrimCol + 1
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
	mGalaxy := galaxy{id: id, distances: make(map[uint]uint32), charges: make(map[uint]uint32), wormholes: make(map[uint16](map[uint]bool))}
	if (id == 0) {
		mGalaxy.distances[0] = 0
	}
	universe.galaxies = append(universe.galaxies, mGalaxy)
}
func addChargeToGalaxy (universe *universe, galaxyId uint, color string, time uint32) {
	id := universe.color2id[color]
	galaxy := universe.galaxies[galaxyId]
	galaxy.charges[id] = time
	/*
	_, hasCharge := galaxy.charges[time]
	if (hasCharge) {
		galaxy.charges[time] = galaxy.charges[time] | id
	} else {
		galaxy.charges[time] = id
	}
	*/
}

func wrapUpGalaxy (universe *universe, galaxyId uint) {
	galaxy := universe.galaxies[galaxyId]
	/*for distance, color := range galaxy.charges {
		for otherDistance, otherColor := range galaxy.charges {
			if (otherDistance < distance) {
				galaxy.charges[distance] = color | otherColor
			}
		}
	}*/
	for _, wormholes := range galaxy.wormholes {
		for color, _:= range wormholes {
			for otherColor, _ := range wormholes {
				if ((otherColor != color) && contained(otherColor, color)) {
					delete(wormholes, color)
				}
			}
		}
	}
}

func addWormHole (universe *universe, color string, start uint16, end uint16) {
	id := universe.color2id[color]
	galaxy := universe.galaxies[start]

	_, isConnected := galaxy.wormholes[end]
	if (!isConnected) {
		galaxy.wormholes[end] = make(map[uint]bool)
		galaxy.wormholes[end][id] = true
	} else {
		galaxy.wormholes[end][id] = true
	}
}
func insertVisit (pq pq.PriorityQueue, visit visit) {
	pq.Insert(visit, visit.distance) // probably should consider color
}

func contained (smallColor uint, bigColor uint)bool {
	return (smallColor & bigColor == smallColor)
}

func formatUniverse (universe universe) []int {
	scores := []int{}
	for _, galaxy := range universe.galaxies {
		scores = append(scores, getGalaxyDistance(galaxy))
	}
	return scores
}

func getGalaxyDistance (galaxy galaxy) int {
	min := -1
	for _, v := range galaxy.distances {
		if (min == -1) {
			min = int(v)
			continue
		} else if (int(v) < min) {
			min = int(v)
		}
	}
	return min
}