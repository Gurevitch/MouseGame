package game

import (
	"fmt"
	"strconv"
	"strings"
)

// InteractionRule is a data-driven NPC reaction: on some trigger, if a
// condition holds, run a list of actions.
//
// Authoring lives in JSON (per-NPC in assets/data/npc/*.json). At runtime,
// evalRules walks a rule list for a given trigger and fires the first one
// whose condition evaluates true.
type InteractionRule struct {
	On     string       `json:"on"`     // trigger: "click" | "item_drop" | "dialog_end" | "state_enter"
	When   string       `json:"when"`   // condition expression; empty = always
	Do     []RuleAction `json:"do"`     // actions to run on match
	Once   bool         `json:"once"`   // fire at most once per NPC-lifetime
	fired  bool         // runtime: whether Once-gated rules already fired
}

// RuleAction is one action in a rule's "do" list. Only the fields relevant
// to the action's Type are populated; the dispatcher ignores the rest.
type RuleAction struct {
	Type   string        `json:"type"`            // dispatch key (see actionDispatch)
	Dialog []dialogEntry `json:"dialog,omitempty"`
	DialogID string      `json:"dialog_id,omitempty"`
	Item   string        `json:"item,omitempty"`
	To     string        `json:"to,omitempty"`
	Scope  string        `json:"scope,omitempty"`
	Key    string        `json:"key,omitempty"`
	Value  int           `json:"value,omitempty"`
	City   string        `json:"city,omitempty"`
	State  string        `json:"state,omitempty"`
	Scene  string        `json:"scene,omitempty"`
	NPC    string        `json:"npc,omitempty"`
	Bool   *bool         `json:"bool,omitempty"`
	Kid    string        `json:"kid,omitempty"`
	Event  string        `json:"event,omitempty"`
	KV     []string      `json:"kv,omitempty"`
}

// ruleContext packages the pieces a condition expression needs to evaluate.
// Passed to evalCondition by fireTrigger.
type ruleContext struct {
	game  *Game
	npc   *npc
	state string // current state machine state; "" if npc has no sm
}

// fireTrigger runs any rule on `n` whose On matches `trigger` and whose
// condition evaluates true. Returns true if any rule fired (useful for
// callers that want to fall back to legacy behavior when no rule matches).
func (g *Game) fireTrigger(n *npc, trigger string, rules []InteractionRule) bool {
	ctx := ruleContext{game: g, npc: n}
	if n != nil && n.sm != nil {
		ctx.state = n.sm.GetState()
	}
	fired := false
	for i := range rules {
		r := &rules[i]
		if r.On != trigger {
			continue
		}
		if r.Once && r.fired {
			continue
		}
		if r.When != "" && !evalCondition(r.When, ctx) {
			continue
		}
		for _, a := range r.Do {
			g.dispatchAction(a, ctx)
		}
		r.fired = true
		fired = true
		// First matching rule wins — more specific rules should be listed first.
		break
	}
	return fired
}

// evalCondition parses and evaluates a condition expression. Supported forms:
//
//	state == <value>
//	inv.has(<item>)
//	chapter.<scope>.<key> == <int>
//	vars.<scope>.<key> == <int>
//	A && B
//	A || B
//
// Unknown expressions return false (safer: unmatched condition = rule skipped).
// Parser is whitespace-tolerant but NOT recursive-descent — it assumes flat
// left-to-right evaluation of && / ||. Good enough for the NPC rules we have;
// if someone needs grouping, rewrite the rule as two simpler rules.
func evalCondition(expr string, ctx ruleContext) bool {
	// OR first (lower precedence)
	if strings.Contains(expr, "||") {
		for _, part := range splitTop(expr, "||") {
			if evalCondition(strings.TrimSpace(part), ctx) {
				return true
			}
		}
		return false
	}
	if strings.Contains(expr, "&&") {
		for _, part := range splitTop(expr, "&&") {
			if !evalCondition(strings.TrimSpace(part), ctx) {
				return false
			}
		}
		return true
	}
	expr = strings.TrimSpace(expr)

	// inv.has(item)
	if strings.HasPrefix(expr, "inv.has(") && strings.HasSuffix(expr, ")") {
		item := strings.TrimSuffix(strings.TrimPrefix(expr, "inv.has("), ")")
		item = strings.TrimSpace(item)
		item = strings.Trim(item, "\"'")
		if ctx.game != nil && ctx.game.inv != nil {
			return ctx.game.inv.hasItem(item)
		}
		return false
	}

	// state == X
	if strings.HasPrefix(expr, "state") {
		left, op, right, ok := splitComparison(expr)
		if !ok || left != "state" {
			return false
		}
		want := strings.Trim(right, "\"' ")
		if op == "==" {
			return ctx.state == want
		}
		if op == "!=" {
			return ctx.state != want
		}
		return false
	}

	// vars.<scope>.<key> == N  or  chapter.<scope>.<key> == N
	if strings.HasPrefix(expr, "vars.") || strings.HasPrefix(expr, "chapter.") {
		left, op, right, ok := splitComparison(expr)
		if !ok {
			return false
		}
		scope, key, ok := splitVarPath(left)
		if !ok {
			return false
		}
		n, err := strconv.Atoi(strings.TrimSpace(right))
		if err != nil || ctx.game == nil || ctx.game.vars == nil {
			return false
		}
		got := ctx.game.vars.Get(scope, key)
		switch op {
		case "==":
			return got == n
		case "!=":
			return got != n
		case ">=":
			return got >= n
		case "<=":
			return got <= n
		case ">":
			return got > n
		case "<":
			return got < n
		}
	}
	return false
}

// splitTop splits `s` on `sep` but only at the top level — it doesn't
// recurse into parens. Used so an expression like `inv.has(x) && state == y`
// splits on && without tearing apart the inv.has(...) call.
func splitTop(s, sep string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '(' {
			depth++
			continue
		}
		if c == ')' {
			depth--
			continue
		}
		if depth == 0 && i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			parts = append(parts, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// splitComparison picks out (left, op, right) from `a <op> b` where op is
// one of ==, !=, <=, >=, <, >. Order matters: check 2-char ops before 1-char.
func splitComparison(expr string) (string, string, string, bool) {
	for _, op := range []string{"==", "!=", "<=", ">=", "<", ">"} {
		if i := strings.Index(expr, op); i >= 0 {
			return strings.TrimSpace(expr[:i]), op, strings.TrimSpace(expr[i+len(op):]), true
		}
	}
	return "", "", "", false
}

// splitVarPath: "vars.game.marcus_healed" -> ("game", "marcus_healed", true)
// Also accepts "chapter.<scope>.<key>" as an alias for "vars.<scope>.<key>".
func splitVarPath(expr string) (string, string, bool) {
	expr = strings.TrimSpace(expr)
	expr = strings.TrimPrefix(expr, "vars.")
	expr = strings.TrimPrefix(expr, "chapter.")
	i := strings.Index(expr, ".")
	if i <= 0 {
		return "", "", false
	}
	return expr[:i], expr[i+1:], true
}

// dispatchAction runs one RuleAction against the game. Unknown types log
// and no-op so a typo in JSON doesn't crash the game.
func (g *Game) dispatchAction(a RuleAction, ctx ruleContext) {
	switch a.Type {
	case "dialog":
		entries := a.Dialog
		if len(entries) == 0 && a.DialogID != "" && g.dialogs != nil {
			// DialogID format is "file.key" (e.g. "marcus.post_strange");
			// dialogStore.Get takes them split.
			if dot := strings.Index(a.DialogID, "."); dot > 0 {
				entries = g.dialogs.Get(a.DialogID[:dot], a.DialogID[dot+1:])
			}
		}
		if len(entries) > 0 {
			g.dialog.startDialog(entries)
		}
	case "queue_dialog":
		if len(a.Dialog) > 0 {
			g.dialog.queueDialog(a.Dialog)
		}
	case "give":
		if a.Item != "" && a.To != "" && g.inv != nil {
			g.inv.giveItemTo(a.Item, a.To)
		}
	case "take":
		if a.Item != "" && g.inv != nil {
			g.inv.removeItem(a.Item)
		}
	case "set_var":
		if a.Scope != "" && a.Key != "" && g.vars != nil {
			g.vars.Set(a.Scope, a.Key, a.Value)
		}
	case "unlock_city":
		if a.City != "" && g.travelMap != nil {
			g.travelMap.setUnlocked(a.City, true)
		}
	case "set_state":
		if ctx.npc != nil && ctx.npc.sm != nil && a.State != "" {
			ctx.npc.sm.SetState(a.State)
		}
	case "set_strange":
		b := a.Bool != nil && *a.Bool
		if ctx.npc != nil {
			ctx.npc.setStrange(b)
		}
	case "set_npc_silent":
		if a.Scene != "" && a.NPC != "" {
			if sc, ok := g.sceneMgr.scenes[a.Scene]; ok {
				for _, n := range sc.npcs {
					if n.name == a.NPC {
						n.silent = a.Bool != nil && *a.Bool
						break
					}
				}
			}
		}
	case "set_scene_npc_strange":
		if a.Scene != "" && a.NPC != "" {
			if sc, ok := g.sceneMgr.scenes[a.Scene]; ok {
				for _, n := range sc.npcs {
					if n.name == a.NPC {
						n.setStrange(a.Bool != nil && *a.Bool)
						break
					}
				}
			}
		}
	case "emit":
		if g.eventBus == nil || a.Event == "" {
			return
		}
		g.eventBus.Emit(EventType(a.Event), a.KV...)
	default:
		fmt.Printf("npc_rules: unknown action %q\n", a.Type)
	}
}
