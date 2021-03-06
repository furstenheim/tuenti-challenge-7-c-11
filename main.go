package main

import "fmt"
import (
	"./priority-queue"
	"strings"
	"strconv"
	"bufio"
	"os"
)


func main() {
	reader := bufio.NewReader(os.Stdin)
	universeLine, _ := reader.ReadString('\n')
	numberOfUniverses, _ := strconv.Atoi(strings.Trim(universeLine, "\n"))
	universeIndex := 1
	for numberOfUniverses > 0 {
		printedIndex := strconv.Itoa(universeIndex)
		os.Stderr.WriteString("--------" + printedIndex + "\n")

		mUniverse := universe{galaxies: []galaxy{}, visits: pq.New(), allColors: make(map[uint]bool), color2id: make(map[string]uint), primary2shift: make(map[string]uint), numPrimCol: 0}
		mUniverse.color2id["Void"] = 0
		mUniverse.allColors[0] = true
		colorsLine, _ := reader.ReadString('\n')
		os.Stderr.WriteString(colorsLine)
		numberOfColours, _ := strconv.Atoi(strings.Trim(colorsLine, "\n"))
		for (numberOfColours > 0) {
			colorLine, _ := reader.ReadString('\n')
			os.Stderr.WriteString(colorLine)
			aColorLine := strings.Split(strings.Trim(colorLine, "\n"), " ")
			colorName := aColorLine[0]
			composed, _ := strconv.Atoi(aColorLine[1])
			colors := []string{}
			if (composed != 0) {
				colors = aColorLine[2: 2 + composed]
			}
			addColorToUniverse(&mUniverse, colorName, uint(composed), colors)
			numberOfColours -= 1
		}
		galaxiesLine, _ := reader.ReadString('\n')
		os.Stderr.WriteString(galaxiesLine)
		numberOfGalaxies, _ := strconv.Atoi(strings.Trim(galaxiesLine, "\n"))
		galaxyId := 0
		for (numberOfGalaxies > 0) {
			createGalaxy(&mUniverse)
			colorsLine, _ := reader.ReadString('\n')
			os.Stderr.WriteString(colorsLine)
			numberOfColours, _ := strconv.Atoi(strings.Trim(colorsLine, "\n"))
			for (numberOfColours > 0) {
				colorLine, _ := reader.ReadString('\n')
				os.Stderr.WriteString(colorLine)
				aColorLine := strings.Split(strings.Trim(colorLine, "\n"), " ")
				colorName := aColorLine[0]
				time, _ := strconv.Atoi(aColorLine[1])
				addChargeToGalaxy(&mUniverse, uint(galaxyId), colorName, uint32(time))
				numberOfColours -=1
			}
			numberOfGalaxies -= 1
			galaxyId += 1
		}
		wormHolesLine, _ := reader.ReadString('\n')
		os.Stderr.WriteString(wormHolesLine)
		numberOfWormholes, _ := strconv.Atoi(strings.Trim(wormHolesLine, "\n"))
		for (numberOfWormholes > 0) {
			wormholeLine, _ := reader.ReadString('\n')
			os.Stderr.WriteString(wormholeLine)
			aWormholeLine := strings.Split(strings.Trim(wormholeLine, "\n"), " ")
			color := aWormholeLine[0]
			idInit, _ := strconv.Atoi(aWormholeLine[1])
			idEnd, _ := strconv.Atoi(aWormholeLine[2])
			addWormHole(&mUniverse, color, uint16(idInit), uint16(idEnd))
			numberOfWormholes -= 1
		}
		for galaxyIndex, _ := range (mUniverse.galaxies) {
			wrapUpGalaxy(&mUniverse, uint(galaxyIndex))
		}
		mUniverse = travelUniverse(mUniverse)
		fmt.Println(formatUniverse(mUniverse, universeIndex))

		numberOfUniverses -= 1
		universeIndex += 1
	}
}

func travelUniverse (universe universe) universe {
	v0 := visit{color: 0, distance: 0, justArrived: true, galaxy: 0}
	insertVisit(universe.visits, v0)
	v1, areThereVisits := universe.visits.Pop()
	for areThereVisits == nil {
		newVisit, _ := v1.(visit)
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
		for color, _ := range universe.allColors {
			if (contained(color, visit.color)) {
				colorDistance, wasVisited := galaxy.distances[color]
				if (!wasVisited || (colorDistance > visit.distance)) {
					galaxy.distances[color] = visit.distance
				}
			}
		}
		return true
	}
	return false
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
	if (!originalVisit.justArrived) {
		return stays
	}
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
			shift, _ := universe.color2id[primary]
			newid = newid | shift
		}
		universe.color2id[name] = newid
	}
	universe.allColors[newid] = true
	for color, _ := range universe.allColors {
		universe.allColors[color | newid] = true
	}
}

func createGalaxy (universe *universe) {
	id := uint16(len(universe.galaxies))
	mGalaxy := galaxy{id: id, distances: make(map[uint]uint32), charges: make(map[uint]uint32), wormholes: make(map[uint16](map[uint]bool))}
	universe.galaxies = append(universe.galaxies, mGalaxy)
}
func addChargeToGalaxy (universe *universe, galaxyId uint, color string, time uint32) {
	id := universe.color2id[color]
	galaxy := universe.galaxies[galaxyId]
	timeToCharge, ok := galaxy.charges[id]
	if (!ok || (timeToCharge > time)) {
		galaxy.charges[id] = time
		// Clean duplicates
		for color, distance := range(galaxy.charges) {
			compounedColor := color | id
			requiredTime := distance + time
			if (contained(color, id)) {
				if (color != id) {
					// Case covered in color === id
					continue
				}
				requiredTime = time
			}
			compounedColorTime, compounedExisted := galaxy.charges[compounedColor]
			if (!compounedExisted || (compounedColorTime > requiredTime)) {
				galaxy.charges[compounedColor] = requiredTime
				for smallColor, smallerDistance :=range (galaxy.charges) {
					if (contained(smallColor, compounedColor) && smallerDistance > requiredTime) {
						delete(galaxy.charges, smallColor)
					}
				}
			}
		}

	}
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

func formatUniverse (universe universe, position int) string {
	scores := []int{}
	for _, galaxy := range universe.galaxies {
		scores = append(scores, getGalaxyDistance(galaxy))
	}
	text := "Case #" + strconv.Itoa(position) + ": " + arrayToString(scores, " ")
	return text
}
func arrayToString(a []int, delim string) string {
    return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
    //return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
    //return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
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
