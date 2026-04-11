package game

// npcState represents a named state for an NPC with dialog, animation, and transitions.
type npcState struct {
	name      string        // state identifier
	dialog    []dialogEntry // what to say in this state
	next      string        // auto-transition to this state after dialog (empty = stay)
	condition func() bool   // only enter if true (nil = always)
}

// npcStateMachine manages NPC dialog states declaratively.
type npcStateMachine struct {
	states   map[string]*npcState
	current  string
	onEnter  map[string]func() // callbacks when entering a state
}

func newNPCStateMachine() *npcStateMachine {
	return &npcStateMachine{
		states:  make(map[string]*npcState),
		onEnter: make(map[string]func()),
	}
}

// AddState registers a named state.
func (sm *npcStateMachine) AddState(name string, dialog []dialogEntry, next string) {
	sm.states[name] = &npcState{
		name:   name,
		dialog: dialog,
		next:   next,
	}
}

// AddStateWithCondition registers a state that only activates when condition is true.
func (sm *npcStateMachine) AddStateWithCondition(name string, dialog []dialogEntry, next string, cond func() bool) {
	sm.states[name] = &npcState{
		name:      name,
		dialog:    dialog,
		next:      next,
		condition: cond,
	}
}

// OnEnter registers a callback for when a state is entered.
func (sm *npcStateMachine) OnEnter(state string, fn func()) {
	sm.onEnter[state] = fn
}

// SetState changes to a new state.
func (sm *npcStateMachine) SetState(name string) {
	sm.current = name
	if fn, ok := sm.onEnter[name]; ok {
		fn()
	}
}

// CurrentDialog returns the dialog for the current state.
func (sm *npcStateMachine) CurrentDialog() []dialogEntry {
	if st, ok := sm.states[sm.current]; ok {
		return st.dialog
	}
	return nil
}

// Advance transitions to the next state (called after dialog ends).
// Returns the name of the new state, or "" if no transition.
func (sm *npcStateMachine) Advance() string {
	st, ok := sm.states[sm.current]
	if !ok || st.next == "" {
		return ""
	}
	nextState, ok := sm.states[st.next]
	if !ok {
		return ""
	}
	if nextState.condition != nil && !nextState.condition() {
		return ""
	}
	sm.SetState(st.next)
	return st.next
}

// GetState returns the current state name.
func (sm *npcStateMachine) GetState() string {
	return sm.current
}
