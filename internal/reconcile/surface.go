package reconcile

import "context"

// Capabilities declares what a surface supports.
type Capabilities struct {
	NoCreate bool // push never creates (e.g. cloud-policies: status reconcile only)
}

// Surface adapts one resource to the engine. All closures are supplied by
// internal/cli; this package never imports an SDK.
type Surface struct {
	Name    string // singular resource noun, e.g. "device rule"
	Command string // CLI group name for guard/audit strings, e.g. "devicecontrol"
	Caps    Capabilities
	// Decode maps one local file into an Object: it unmarshals into the
	// surface's typed file-shape struct, validates the identity fields are
	// present (hard error otherwise), and re-marshals the struct as the
	// canonical Body. List produces live Objects through the SAME file-shape
	// struct + marshal, so bodies are byte-comparable.
	Decode func(data []byte) (Object, error)
	List   func(ctx context.Context) ([]Object, error)   // live objects
	Create func(ctx context.Context, local Object) error // nil if NoCreate
	Update func(ctx context.Context, id string, local Object) error
}
