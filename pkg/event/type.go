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
)

const (
	internalEvtTypeTerm = "internal:evt:term"
)

type EventType string

// CallbackFn is the function that can be registered to be automatically
// called during the handling of a given type of event. All callbacks
// have a context, a pointer to the event that triggered the callback
// and a slice of slice of bytes to pass in any data.
type CallbackFn func(context.Context, *Engine, *Event) error

func (e *Engine) RegisterCallback(t *EventType, cb CallbackFn, args ...interface{}) error {
	if !e.typeExists(*t) {
		return fmt.Errorf("unknown event type")
	}

	e.types[*t] = append(e.types[*t], cb)
	return nil
}

func (e *Engine) typeExists(typeID EventType) bool {
	if _, ok := e.types[typeID]; ok {
		return true
	}
	return false
}

func (e *Engine) NewType(id string) (EventType, error) {
	var et EventType
	et = EventType(id)
	if e.typeExists(et) {
		return et, fmt.Errorf("type %s already exists", id)
	}
	e.types[et] = nil
	return et, nil
}

func (evt *Event) SetType(t EventType) error {
	if evt == nil {
		return fmt.Errorf("undefined event")
	}
	evt.EventType = t
	return nil
}
