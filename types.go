package main

import "./priority-queue"

type galaxy struct {
	id uint16
	distances map[uint]uint32
	charges map[uint]uint32
	wormholes map[uint16](map[uint]bool) // galaxy and number of possible colors
}

type visit struct {
	color uint
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
