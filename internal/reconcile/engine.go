package reconcile

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Summary is the apply outcome.
type Summary struct{ Created, Updated, Failed int }

// Apply executes the plan's creates and updates with warn-and-continue
// semantics, writing per-item failures to warnw. Creates are processed before
// updates, each sorted by Name for deterministic order. A create on a NoCreate
// surface (Create == nil) counts as a failure ("cannot create <noun> <name>").
// Returns a non-nil error when Failed > 0 ("N of M changes failed") so callers
// exit non-zero and guard() audits the failure.
func Apply(ctx context.Context, s Surface, p Plan, warnw io.Writer) (Summary, error) {
	creates := p.Creates()
	updates := p.Updates()
	sortByName(creates)
	sortByName(updates)

	var sum Summary
	for _, it := range creates {
		if it.Local == nil {
			fmt.Fprintf(warnw, "cannot create %s %q: missing local definition\n", s.Name, it.Name)
			sum.Failed++
			continue
		}
		if s.Create == nil {
			fmt.Fprintf(warnw, "cannot create %s %q: surface does not support create\n", s.Name, it.Name)
			sum.Failed++
			continue
		}
		if err := s.Create(ctx, *it.Local); err != nil {
			fmt.Fprintf(warnw, "failed to create %s %q: %v\n", s.Name, it.Name, err)
			sum.Failed++
			continue
		}
		sum.Created++
	}
	for _, it := range updates {
		if it.Local == nil {
			fmt.Fprintf(warnw, "cannot update %s %q: missing local definition\n", s.Name, it.Name)
			sum.Failed++
			continue
		}
		if err := s.Update(ctx, it.ID, *it.Local); err != nil {
			fmt.Fprintf(warnw, "failed to update %s %q: %v\n", s.Name, it.Name, err)
			sum.Failed++
			continue
		}
		sum.Updated++
	}

	if sum.Failed > 0 {
		total := len(creates) + len(updates)
		return sum, fmt.Errorf("%d of %d changes failed", sum.Failed, total)
	}
	return sum, nil
}

func sortByName(items []Item) {
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
}

// LoadDir reads every *.yaml/*.yml file in dir into Objects via decode. A file
// that fails to decode is a hard error naming the file. A missing dir returns
// (nil, nil).
func LoadDir(dir string, decode func(data []byte) (Object, error)) ([]Object, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var objs []Object
	for _, e := range entries {
		if e.IsDir() || !isYAML(e.Name()) {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, rErr := os.ReadFile(path)
		if rErr != nil {
			return nil, rErr
		}
		o, dErr := decode(data)
		if dErr != nil {
			return nil, fmt.Errorf("decode %s: %w", path, dErr)
		}
		objs = append(objs, o)
	}
	return objs, nil
}

// WriteDir writes each object to dir as <sanitizeName(o.Name)>.yaml with
// duplicate stems suffixed -1, -2… (same behavior as the legacy pulls). It
// returns the names of pre-existing *.yaml/*.yml files in dir that did NOT
// correspond to a written object (stale files) so pull can warn: a stale file
// would plan as a create on the next push.
func WriteDir(dir string, objects []Object) (stale []string, err error) {
	existing := make(map[string]bool)
	if entries, rErr := os.ReadDir(dir); rErr == nil {
		for _, e := range entries {
			if !e.IsDir() && isYAML(e.Name()) {
				existing[e.Name()] = true
			}
		}
	} else if !os.IsNotExist(rErr) {
		return nil, rErr
	}

	if mErr := os.MkdirAll(dir, 0o750); mErr != nil {
		return nil, mErr
	}

	used := make(map[string]int)
	written := make(map[string]bool)
	for _, o := range objects {
		base := sanitizeName(o.Name)
		stem := base
		if n := used[base]; n > 0 {
			stem = fmt.Sprintf("%s-%d", base, n)
		}
		used[base]++

		name := stem + ".yaml"
		if wErr := os.WriteFile(filepath.Join(dir, name), o.Body, 0o644); wErr != nil {
			return nil, wErr
		}
		written[name] = true
	}

	for name := range existing {
		if !written[name] {
			stale = append(stale, name)
		}
	}
	sort.Strings(stale)
	return stale, nil
}

func isYAML(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".yaml" || ext == ".yml"
}

var unsafeNameChars = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

// sanitizeName converts an object name into a safe filename stem, mirroring the
// legacy per-object pull behavior.
func sanitizeName(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = unsafeNameChars.ReplaceAllString(s, "")
	if s == "" {
		s = "object"
	}
	return s
}
