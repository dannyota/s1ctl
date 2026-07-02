package reconcile

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Item kinds.
const (
	KindCreate    = "create"
	KindUpdate    = "update"
	KindUnchanged = "unchanged"
	KindLiveOnly  = "live-only"
)

// Item is one planned change. Kind is one of the Kind* constants.
type Item struct {
	Kind  string
	Name  string
	ID    string
	Local *Object
	Live  *Object
}

// Plan is the classified set of items produced by BuildPlan.
type Plan struct{ Items []Item }

// BuildPlan classifies local files against live objects, matching by
// Object.Name. Byte-equal bodies → unchanged. Local-only → create. Live-only →
// reported, never deleted. Duplicate live names: first wins (returned in the
// warnings slice). Duplicate LOCAL names: hard error (ambiguous input, nothing
// applied). Input order does not affect the result: items are sorted by name.
func BuildPlan(local, live []Object) (Plan, []string, error) {
	localByName := make(map[string]Object, len(local))
	localNames := make([]string, 0, len(local))
	for _, o := range local {
		if _, dup := localByName[o.Name]; dup {
			return Plan{}, nil, fmt.Errorf("duplicate local name %q: input is ambiguous, nothing applied", o.Name)
		}
		localByName[o.Name] = o
		localNames = append(localNames, o.Name)
	}

	liveByName := make(map[string]Object, len(live))
	liveNames := make([]string, 0, len(live))
	var warnings []string
	for _, o := range live {
		if _, seen := liveByName[o.Name]; seen {
			warnings = append(warnings, fmt.Sprintf("duplicate live name %q: keeping first, ignoring the rest", o.Name))
			continue
		}
		liveByName[o.Name] = o
		liveNames = append(liveNames, o.Name)
	}

	sort.Strings(localNames)
	sort.Strings(liveNames)

	var items []Item
	for _, name := range localNames {
		l := localByName[name]
		lv, ok := liveByName[name]
		if !ok {
			lc := l
			items = append(items, Item{Kind: KindCreate, Name: name, ID: l.ID, Local: &lc})
			continue
		}
		lc, vc := l, lv
		if bytes.Equal(l.Body, lv.Body) {
			items = append(items, Item{Kind: KindUnchanged, Name: name, ID: lv.ID, Local: &lc, Live: &vc})
			continue
		}
		items = append(items, Item{Kind: KindUpdate, Name: name, ID: lv.ID, Local: &lc, Live: &vc})
	}
	for _, name := range liveNames {
		if _, ok := localByName[name]; ok {
			continue
		}
		vc := liveByName[name]
		items = append(items, Item{Kind: KindLiveOnly, Name: name, ID: vc.ID, Live: &vc})
	}

	return Plan{Items: items}, warnings, nil
}

func (p Plan) filter(kind string) []Item {
	var out []Item
	for _, it := range p.Items {
		if it.Kind == kind {
			out = append(out, it)
		}
	}
	return out
}

// Creates returns the create items.
func (p Plan) Creates() []Item { return p.filter(KindCreate) }

// Updates returns the update items.
func (p Plan) Updates() []Item { return p.filter(KindUpdate) }

// Unchanged returns the unchanged items.
func (p Plan) Unchanged() []Item { return p.filter(KindUnchanged) }

// LiveOnly returns the live-only items (reported, never deleted).
func (p Plan) LiveOnly() []Item { return p.filter(KindLiveOnly) }

// Empty reports whether the plan has no creates and no updates.
func (p Plan) Empty() bool {
	for _, it := range p.Items {
		if it.Kind == KindCreate || it.Kind == KindUpdate {
			return false
		}
	}
	return true
}

// Describe renders the guard action string for a plan, e.g.
// "create 2 device rules, update 1 device rule (3 unchanged) from <dir>".
// Zero-count create/update segments are omitted; when neither is present it
// reads "no changes for <nouns> in <dir>".
func Describe(p Plan, noun, dir string) string {
	creates := len(p.Creates())
	updates := len(p.Updates())
	unchanged := len(p.Unchanged())

	if creates == 0 && updates == 0 {
		return fmt.Sprintf("no changes for %s in %s", nounForm(2, noun), dir)
	}

	var segs []string
	if creates > 0 {
		segs = append(segs, "create "+pluralize(creates, noun))
	}
	if updates > 0 {
		segs = append(segs, "update "+pluralize(updates, noun))
	}
	s := strings.Join(segs, ", ")
	if unchanged > 0 {
		s += fmt.Sprintf(" (%d unchanged)", unchanged)
	}
	return s + " from " + dir
}

// pluralize renders a count with its naively pluralized noun, e.g.
// pluralize(2, "device rule") → "2 device rules".
func pluralize(n int, singular string) string {
	return fmt.Sprintf("%d %s", n, nounForm(n, singular))
}

// nounForm returns singular for n == 1 and a naive plural otherwise.
func nounForm(n int, singular string) string {
	if n == 1 {
		return singular
	}
	if last := singular[len(singular)-1]; last == 'y' && len(singular) > 1 {
		prev := singular[len(singular)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return singular[:len(singular)-1] + "ies"
		}
	}
	return singular + "s"
}
