// Package reconcile is the SDK-free core behind every config-as-code surface.
//
// It owns matching, planning, file I/O, and apply semantics for the core loop:
// pull live state, review the diff in git, push back. Each surface supplies
// closures that adapt one resource to the engine; this package imports no SDK
// and performs no HTTP, so it is fully unit-testable with fake closures.
package reconcile

// Object is one config item in canonical file form.
type Object struct {
	Name string // stable identity used for matching (surface-defined)
	ID   string // server ID; "" for local objects not yet created
	Body []byte // canonical declarative body (YAML of the surface's file shape)
}
