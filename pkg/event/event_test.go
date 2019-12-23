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
	"testing"
)

var events []*Event

func callback(ctx context.Context, engine *Engine, evt *Event) error {
	fmt.Printf("Returning event %d\n", evt.ID)

	err := engine.Return(evt)
	if err != nil {
		return fmt.Errorf("unable to return event: %s", err)
	}

	return nil
}

func callbackWithTracking(ctx context.Context, engine *Engine, evt *Event) error {
	idx := evt.ID
	err := callback(ctx, engine, evt)
	if err != nil {
		return err
	}
	events[idx] = nil

	return nil
}

func testCompleted(evts []*Event) bool {
	for i := 0; i < len(evts); i++ {
		if evts[i] != nil {
			return false
		}
	}

	return true
}

func initEngine(t *testing.T, queueCfg QueueCfg) (*Engine, EventType) {

	log.Println("Creating event engine...")
	engine := queueCfg.Init()
	if engine == nil {
		t.Fatal("unable to create event engine")
	}

	// Create a new type of events
	log.Println("Creating new event type...")
	eventType, err := engine.NewType("dummy")
	if err != nil {
		t.Fatal("unable to create new event type: %w", err)
	}

	return engine, eventType
}

func TestConsumeReturnEvents(t *testing.T) {
	// Create a new engine with an initial queue of 1024 inactive events
	queueCfg := QueueCfg{
		Size: 1024,
	}

	engine, eventType := initEngine(t, queueCfg)
	if engine == nil {
		t.Fatal("unable to initialize engine")
	}

	log.Println("Registering callback...")
	err := engine.RegisterCallback(&eventType, callbackWithTracking)
	if err != nil {
		t.Fatal("unable to register callback: %w", err)
	}

	var i uint64
	log.Println("Creating and emitting events...")
	events = make([]*Event, queueCfg.Size)
	for i = 0; i < queueCfg.Size; i++ {
		events[i] = engine.GetEvent(true)
		err := events[i].SetType(eventType)
		if err != nil {
			t.Fatalf("failed to set event type")
		}
		err = events[i].Emit(nil)
		if err != nil {
			t.Fatalf("failed to emit event")
		}
	}

	log.Println("All done terminating")
	for {
		if testCompleted(events) {
			break
		}
	}

	engine.Fini()
}

func TestEventExhaustion(t *testing.T) {
	queueCfg := QueueCfg{
		Size: 1,
	}

	engine, eventType := initEngine(t, queueCfg)
	if engine == nil {
		t.Fatal("unable to initialize engine")
	}

	log.Println("Registering callback...")
	err := engine.RegisterCallback(&eventType, callback)
	if err != nil {
		t.Fatal("unable to register callback: %w", err)
	}

	t.Log("Getting first event, expected to succeed...")
	evt1 := engine.GetEvent(true)
	if evt1 == nil {
		t.Fatal("unable to get first event")
	}
	err = evt1.SetType(eventType)
	if err != nil {
		t.Fatalf("failed to set event type")
	}
	t.Log("Getting second event, expected to fail...")
	evt2 := engine.GetEvent(false)
	if evt2 != nil {
		t.Fatal("we were able to get an event while we were expecting to have no free event")
	}
	err = evt1.Emit(nil)
	if err != nil {
		t.Fatal("unable to emit event: %w", err)
	}
	evt2 = engine.GetEvent(true)
	if evt2 == nil {
		t.Fatal("unable to get the second event")
	}
	err = evt2.SetType(eventType)
	if err != nil {
		t.Fatalf("failed to set event type")
	}
	err = evt2.Emit(nil)
	if err != nil {
		t.Fatal("unable to emit event: %w", err)
	}
	engine.Fini()
}
