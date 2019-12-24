/*
 * Copyright (c) 2019 Geoffroy Vallee, All rights reserved
 * This software is licensed under a 3-clause BSD license. Please consult the
 * LICENSE.md file distributed with the sources of this project regarding your
 * rights to use or distribute this software.
 */

package event

import (
	"fmt"
)

// MAXARGS is the maximim number of arguments that can be associated to
// any event. Each argument is a slice of bytes of any length
const MAXARGS = 8

// Event is a structure representing an event
type Event struct {
	ID        uint64          //`json:"id"`
	EventType EventType       //`json:"event_type"`
	Data      [MAXARGS][]byte //`json:"data"`
	engine    *Engine
}

// Emit triggers an event, i.e., the event will be added to the list of active
// event until it is dispatched by the engine. When being dispatched, the
// callbacks registered to the event type are automatically called.
func (evt *Event) Emit(data []byte) error {
	if evt == nil {
		return fmt.Errorf("undefined event")
	}

	evt.engine.activeQueue <- *evt

	return nil
}
