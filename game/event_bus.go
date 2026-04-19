package game

// EventBus is a typed pub/sub hub used to decouple subsystems from story
// progression. Publishers (NPCs, sequences, scene loader) emit events;
// subscribers (chapter wiring in paris.go / jerusalem.go / ...) react.
//
// It is intentionally dependency-free and synchronous: handlers run in the
// same goroutine that publishes the event, in registration order. No queues,
// no threading — the game loop is single-threaded and sequences expect
// deterministic ordering.
//
// Status: infrastructure-only in Phase 1. Not wired into existing code yet.
// Phases 4+ replace the closure tree in setupCampCallbacks with subscriptions
// on this bus.

type EventType string

const (
	EvtDialogEnded  EventType = "dialog_ended"   // payload: NPCID
	EvtItemGiven    EventType = "item_given"     // payload: Item, To
	EvtCityUnlocked EventType = "city_unlocked"  // payload: City
	EvtKidHealed    EventType = "kid_healed"     // payload: Kid
	EvtSceneEntered EventType = "scene_entered"  // payload: Scene
	EvtSceneExited  EventType = "scene_exited"   // payload: Scene
	EvtChapterStart EventType = "chapter_start"  // payload: Chapter
	EvtChapterEnd   EventType = "chapter_end"    // payload: Chapter
)

// Event carries a type and a string-keyed payload. Keeping the payload as a
// map (rather than a typed struct per event) lets handlers ignore fields they
// do not care about and lets JSON-authored rules emit events without Go code.
type Event struct {
	Type    EventType
	Payload map[string]string
}

// Handler is a function that reacts to one event.
type Handler func(Event)

// EventBus is the hub.
type EventBus struct {
	subs map[EventType][]Handler
}

func newEventBus() *EventBus {
	return &EventBus{subs: make(map[EventType][]Handler)}
}

// Subscribe registers a handler for an event type. Returns an unsubscribe
// function; callers that own transient subscriptions (e.g. a one-shot scene
// listener) should call it when their scope ends.
func (b *EventBus) Subscribe(t EventType, h Handler) func() {
	b.subs[t] = append(b.subs[t], h)
	idx := len(b.subs[t]) - 1
	return func() {
		// Best-effort unsubscribe: nil the slot so the handler stops firing.
		// We do not compact the slice because handlers are few and long-lived.
		if idx < len(b.subs[t]) {
			b.subs[t][idx] = nil
		}
	}
}

// Publish dispatches an event synchronously to all current subscribers.
// Handlers registered during dispatch fire on the *next* publish of the same
// type, not this one — preventing re-entrant surprises.
func (b *EventBus) Publish(t EventType, payload map[string]string) {
	handlers := b.subs[t]
	evt := Event{Type: t, Payload: payload}
	for _, h := range handlers {
		if h != nil {
			h(evt)
		}
	}
}

// Emit is a convenience for the common "type + key/value pairs" form.
// Odd number of args is ignored silently so rule authors can't crash the game
// with a typo.
func (b *EventBus) Emit(t EventType, kv ...string) {
	payload := map[string]string{}
	for i := 0; i+1 < len(kv); i += 2 {
		payload[kv[i]] = kv[i+1]
	}
	b.Publish(t, payload)
}
