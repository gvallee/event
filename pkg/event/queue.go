/*
 * Copyright (c) 2019 Geoffroy Vallee, All rights reserved
 * This software is licensed under a 3-clause BSD license. Please consult the
 * LICENSE.md file distributed with the sources of this project regarding your
 * rights to use or distribute this software.
 */

package event

import "fmt"

type Queue chan Event

type QueueCfg struct {
	Size uint64
}

// Init creates a new queue of event with n initial events in it
func InitQueue(cfg QueueCfg) Queue {
	q := make(chan Event, cfg.Size)
	var i uint64
	for i = 0; i < cfg.Size; i++ {
		var e Event
		e.ID = i
		q <- e
	}

	return q
}

// Pull is a non-blocking function that tries to get an event from a queue
func Pull(q *Queue) *Event {
	if q == nil {
		return nil
	}

	if len(*q) == 0 {
		return nil
	}

	e := <-(*q)
	return &e
}

// Pool will block until an event is available on the queue
func Poll(q *Queue) *Event {
	if q == nil {
		return nil
	}
	e := <-(*q)
	return &e
}

// Return appends an event to a queue
func (q *Queue) Return(e *Event) error {
	if q == nil {
		return fmt.Errorf("undefined queue")
	}

	if e == nil {
		return nil
	}

	*q <- *e
	e = nil

	return nil
}
