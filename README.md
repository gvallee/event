# event
A Go module providing an event system

# Overview

This module provides a fairly classic event system.
The core concepts are:
- An engine: each engine has two event queues; one for events that can be handed
to the calling application for usage; one of active event, that needs for the 
engine to be handled.
- For each given engine, one can register a new event type. Event types are 
designed to guve the opportunity to applications to attach semantics to a given
event.
- For each event type, one can register callbacks. If multiple callbacks are 
registered, they are called in order of the registration, it is the 
responsability of the calling applications to know that order when relevant (by
controling how callbacks are registered). If a callback returns an error, the 
engine will stop.
- The calling application gets an event from the engine before being able to 
use it. Two modes are available, blocking and non-blocking. In blocking mode, 
the call to get the event blocks until an event is available; while in 
non-blocking, the call will return nil if no event is available.
- Once the application is done with an event, the event must be returned so it
can be reused.

# Usage

For more details about usage, please refer to tests.

## Create a new engine

```
queueCfg := QueueCfg{
	Size: <initial number of available events>,
}
engine := queueCfg.Init()
if engine == nil {
	return fmt.Errorf("unable to create event engine")
}
```

## Termination of an engine

```
engine.Fini()
```

## Create of a new event type and registration of a callback

```
func callback(ctx context.Context, engine *Engine, evt *Event) error {
    // A simple example of a callback returning an event
	err := engine.Return(evt)
	if err != nil {
		return fmt.Errorf("unable to return event: %w", err)
	}
	return nil
}
```
```
eventType, err := engine.NewType("myEventType")
if err != nil {
	return fmt.Errorf("unable to create new event type: %w", err)
}
err := engine.RegisterCallback(&eventType, callback)
if err != nil {
	return fmt.Errorf("unable to register callback: %w", err)
}
```

# Get an event and emit it

```
evt1 := engine.GetEvent(true)
if evt1 == nil {
	return fmt.Errorf("unable to get first event")
}
err = evt1.SetType(eventType)
if err != nil {
	return fmt.Errorf("failed to set event type")
}
err = evt1.Emit(nil)
if err != nil {
	return fmt.Errorf("unable to emit event: %w", err)
}
```

