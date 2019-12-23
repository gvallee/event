/*
 * Copyright (c) 2019 Geoffroy Vallee, All rights reserved
 * This software is licensed under a 3-clause BSD license. Please consult the
 * LICENSE.md file distributed with the sources of this project regarding your
 * rights to use or distribute this software.
 */

package event

import (
	"context"
	"fmt"
	"log"
)

type Engine struct {
	InactiveEvts Queue
	activeQueue  Queue
	types        map[EventType][]CallbackFn
}

func (e *Engine) handleActiveEvents() {
	for {
		log.Println("Waiting for an active event")
		evt := <-e.activeQueue
		log.Println("Active event available")

		if !e.typeExists(evt.EventType) {
			return
		}

		callbacks := e.types[evt.EventType]
		for _, cb := range callbacks {
			var ctx context.Context

			if evt.EventType == internalEvtTypeTerm {
				log.Println("Termination event received, returning...")
				return
			}

			fmt.Println("Calling event's callback...")
			err := cb(ctx, e, &evt)
			if err != nil {
				log.Println("callback failed: %w", err)
			}
		}
	}
}

func (cfg QueueCfg) Init() *Engine {
	var e Engine

	log.Println("Initializing initial events...")
	e.InactiveEvts = InitQueue(cfg)
	if e.InactiveEvts == nil {
		return nil
	}

	// activeQueueCfg is the configuration of the queue that is
	// used to handle event that has been emitted but not dispatched yet
	log.Println("Initializing the queue of active events...")
	activeQueueCfg := QueueCfg{
		Size: 0,
	}
	e.activeQueue = InitQueue(activeQueueCfg)

	// Add the default event types
	e.types = make(map[EventType][]CallbackFn)
	_, err := e.NewType(internalEvtTypeTerm)
	if err != nil {
		return nil
	}

	// Start a go routine that will block until events are available in the
	// activeQueue
	log.Println("Creating the thread to handle events...")
	go e.handleActiveEvents()

	log.Println("Initialization succeeded")
	return &e
}

func (e *Engine) Return(evt *Event) error {
	return e.InactiveEvts.Return(evt)
}

func (e *Engine) Fini() {
	evt := <-e.InactiveEvts
	err := evt.SetType(internalEvtTypeTerm)
	if err != nil {
		log.Println("[ERROR] failed to set the type for a termination event")
	}
	err = evt.Emit(nil)
	if err != nil {
		log.Println("[ERROR] failed to emit termination event")
	}
}

func (e *Engine) GetEvent(block bool) *Event {
	if e == nil {
		return nil
	}

	var evt *Event
	if block {
		evt = Poll(&e.InactiveEvts)
		evt.engine = e
	} else {
		evt = Pull(&e.InactiveEvts)
		if evt != nil {
			evt.engine = e
		}
	}

	return evt
}
