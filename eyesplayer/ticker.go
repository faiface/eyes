package main

import "time"

type UniformTicker struct {
	Step     time.Duration
	Action   func()
	time     time.Duration
	nextTick time.Duration
}

func (ut *UniformTicker) Update(dt time.Duration) {
	ut.time += dt
	for ut.nextTick < ut.time {
		ut.Action()
		ut.nextTick += ut.Step
	}
}
